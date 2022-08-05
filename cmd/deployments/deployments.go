package deployments

import (
	"fmt"
	"log"
	"reflect"
	"time"

	od "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	p "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
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
octostats deployments --project 'Etrm.Til.FileSystemConnector' --environment 'Etrm Production'

# Get deployment stats for all projects in a project group
octostats deployments --projectgroup 'Etrm.Integration' --environment 'Etrm Production'`
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
			// fmt.Println(contains([]*p.Project{{Name: "bla"}}, &p.Project{Name: "bla"}))
		},
	}

	cmd.Flags().StringVarP(&o.Project, "project", "p", "", "Name of the Project")
	cmd.Flags().StringVarP(&o.ProjectGroup, "projectgroup", "g", "", "Name of the Project Group")
	cmd.Flags().StringVarP(&o.Environment, "environment", "e", "", "Environment for which to gather statistics")
	cmd.Flags().IntVar(&o.Lookback, "lookback", 30, "How many days to look back for deployments")

	return cmd
}

func (o *DeploymentOptions) init() {
	if o.Environment != "" { // only calculate the environment is not empty
		environmentId = environment.GetEnvironmentID(o.Client, o.Environment)
	}
}

func (o *DeploymentOptions) ShowDeploymentStats() {
	var projects []*p.Project

	if o.ProjectGroup != "" {
		pg, err := util.GetProjectGroupByName(*o.Client, o.ProjectGroup)
		if err != nil {
			log.Fatalln(fmt.Errorf("error getting project group: %s", o.ProjectGroup))
		}

		ps, err := o.Client.ProjectGroups.GetProjects(pg)
		if err != nil {
			log.Fatalln(fmt.Errorf("error getting project in project group: %s", pg.Name))
		}
		projects = append(projects, ps...)
	}

	if o.Project != "" {
		project, err := util.GetProjectByName(o.Client, o.Project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching project: %s", o.Project))
		}
		if !contains(projects, project) {
			projects = append(projects, project)
		}
	}

	count := 0
	for _, project := range projects {
		pCount := 0
		releases, err := o.Client.Projects.GetReleases(project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching releases for project: %s", o.Project))
		}

		for i := len(releases) - 1; i >= 0; i-- { // start with the earliest release
			release := releases[i]
			deployments, err := o.Client.Deployments.GetDeployments(release)
			if err != nil {
				log.Fatalln(fmt.Errorf("error getting deployments for release: %s", release.ID))
			}
			for _, d := range deployments.Items {
				if environmentId == "" || environmentId == d.EnvironmentID {
					if d.Created.After(time.Now().AddDate(0, 0, 0-o.Lookback)) {
						if pCount == 0 {
							fmt.Println(project.Name) // print project name on first valid release found
						}
						fmt.Printf("	%v: %v\n", d.Created.Format("02/01"), release.Version)
						pCount++
						break
					}
				}
			}
		}
		count += pCount
	}

	fmt.Printf("Number of releases in the past %v days: %v\n", o.Lookback, count)
}

func contains[T comparable](elems []T, v T) bool {
	// fmt.Println(v)
	for _, s := range elems {
		// fmt.Println(s)
		if reflect.DeepEqual(s, v) {
			return true
		}
	}
	return false
}
