package deployments

import (
	"fmt"
	"log"
	"time"

	od "github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/spf13/cobra"
	"n8m.io/octostats/pkg/environment"
	"n8m.io/octostats/pkg/util"
)

type DeploymentOptions struct {
	Project      string
	ProjectGroup string
	Environment  string
	Lookback     int

	Client *od.Client
}

var (
	deploymentsExample = `
# Get deployment stats for a project
octostats deployments --project 'Etrm.Til.FileSystemConnector --environment "Etrm Production"'

# Get deployment stats for all projects in a project group
octostats deployments --projectgroup 'Etrm.Integration'`
	environmentId string
)

func NewDeploymentsCmd() *cobra.Command {
	o := &DeploymentOptions{}

	var cmd = &cobra.Command{
		Use:     "deployments",
		Short:   "Get deployment frequency statistics",
		Example: deploymentsExample,
		PreRun: func(cmd *cobra.Command, args []string) {
			o.Client = util.OctopusClient()
			o.init()
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.ShowDeploymentStats()
		},
	}

	cmd.Flags().StringVarP(&o.Project, "project", "p", "", "Name of the Project")
	cmd.Flags().StringVarP(&o.ProjectGroup, "projectgroup", "g", "", "Name of the Project Group")
	cmd.Flags().StringVarP(&o.Environment, "environment", "e", "", "Environment for which to gather statistics")
	cmd.Flags().IntVar(&o.Lookback, "lookback", 30, "How many days to look back for deployments")

	return cmd
}

func (o *DeploymentOptions) init() {
	ids, err := environment.GetEnvironmentIDs(o.Client, o.Environment)
	if err != nil {
		log.Fatalln(fmt.Errorf("error fetching ID for environment: %s", o.Environment))
	} else if len(ids) > 1 {
		log.Fatalln(fmt.Errorf("found more than one environment matching: %s", o.Environment))
	}
	environmentId = ids[0]
}

func (o *DeploymentOptions) ShowDeploymentStats() {
	var releases []*od.Release
	if o.Project != "" {
		project, err := util.GetProjectByName(o.Client, o.Project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching project: %s", o.Project))
		}

		releases, err = o.Client.Projects.GetReleases(project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching releases for project: %s", o.Project))
		}

	} else if o.ProjectGroup != "" {
		// projectGroups, err := o.Client.ProjectGroups.GetByPartialName(o.ProjectGroup)
		// if err != nil {
		// 	log.Fatalln(fmt.Errorf("error fetching project group: %s", o.ProjectGroup))
		// }

		// TODO: Get all releases in project group, no easy way so far
	}

	count := 0
	for _, r := range releases {
		deployments, err := o.Client.Deployments.GetDeployments(r)
		if err != nil {
			log.Fatalln(fmt.Errorf("error getting deployments for release: %s", r.ID))
		}
		for _, d := range deployments.Items {
			if environmentId == d.EnvironmentID && d.Created.After(time.Now().AddDate(0, 0, 0-o.Lookback)) {
				count++
				break
			}
		}
	}

	fmt.Printf("Number of releases in the past %v days: %v\n", o.Lookback, count)
}

// func (o *DeploymentOptions) GetProjectGroupReleases() {
// 	projectGroups, _ := o.Client.ProjectGroups.Get(od.ProjectGroupsQuery{PartialName: o.ProjectGroup})
// 	if len(projectGroups.Items) > 1 {
// 		log.Fatalln(fmt.Errorf("got more than one project group matching pattern"))
// 	}

// 	query := od.EventsQuery{
// 		ProjectGroups:   []string{projectGroups.Items[0].ID},
// 		EventCategories: []string{"DeploymentStarted"},
// 	}

// 	events, err := o.Client.Events.Get(query)
// 	if err != nil {
// 		log.Fatalln(fmt.Errorf("error getting events for project group: %s", o.ProjectGroup))
// 	}
// 	if len(events.Items) == 0 {
// 		fmt.Println("no events")
// 	}
// 	for _, e := range events.Items {
// 		fmt.Printf("%v, %v\n", e.Category, e.Occurred.Format("02-01-2006"))
// 		e.Details
// 	}
// }
