package deployments

import (
	"fmt"
	"log"
	"time"

	od "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/releases"
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
	environmentId = environment.GetEnvironmentID(o.Client, o.Environment)
}

func (o *DeploymentOptions) ShowDeploymentStats() {
	var releases []*releases.Release
	var projects []*projects.Project
	if o.Project != "" {
		project, err := util.GetProjectByName(o.Client, o.Project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching project: %s", o.Project))
		}
		projects = append(projects, project)

	} else if o.ProjectGroup != "" {
		projectGroups, err := o.Client.ProjectGroups.GetByPartialName(o.ProjectGroup)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching project group: %s", o.ProjectGroup))
		} else if len(projectGroups) > 1 {
			log.Fatalln(fmt.Errorf("found more than one project group matching: %s", o.Environment))
		}
		pg := projectGroups[0]

		ps, err := o.Client.ProjectGroups.GetProjects(pg)
		if err != nil {
			log.Fatalln(fmt.Errorf("error getting project in project group: %s", pg.Name))
		}
		projects = append(projects, ps...)
	}

	for _, project := range projects {
		fmt.Println("Looking at project: " + project.Name)
		rs, err := o.Client.Projects.GetReleases(project)
		if err != nil {
			log.Fatalln(fmt.Errorf("error fetching releases for project: %s", o.Project))
		}
		releases = append(releases, rs...)
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
