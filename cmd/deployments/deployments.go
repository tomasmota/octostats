package deployments

import (
	"fmt"
	"log"
	"time"

	od "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
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
	var projects []*projects.Project
	if o.Project != "" {
		project, err := util.GetProjectByName(o.Client, o.Project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching project: %s", o.Project))
		}
		projects = append(projects, project)

	} else if o.ProjectGroup != "" {
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

	count := 0
	for _, project := range projects {
		fmt.Printf("project: %v\n", project.Name)
		releases, err := o.Client.Projects.GetReleases(project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching releases for project: %s", o.Project))
		}

		for _, release := range releases {
			deployments, err := o.Client.Deployments.GetDeployments(release)
			if err != nil {
				log.Fatalln(fmt.Errorf("error getting deployments for release: %s", release.ID))
			}
			for _, d := range deployments.Items {
				if environmentId == "" || environmentId == d.EnvironmentID {
					if d.Created.After(time.Now().AddDate(0, 0, 0-o.Lookback)) {
						fmt.Println("  " + release.Version)
						count++
						break
					}
				}
			}
		}
	}

	if o.Environment != "" {
		fmt.Printf("\nNumber of releases to '%v' in the past %v days: %v\n", o.Environment, o.Lookback, count)
	} else {
		fmt.Printf("\nNumber of releases across all environments in the past %v days: %v\n", o.Lookback, count)
	}
}
