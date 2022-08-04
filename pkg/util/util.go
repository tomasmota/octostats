package util

import (
	"fmt"
	"log"
	"net/url"

	od "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/spf13/viper"
)

func GetProjectGroupByName(client od.Client, name string) (*projectgroups.ProjectGroup, error) {
	projectGroups, err := client.ProjectGroups.GetByPartialName(name)
	if err != nil {
		log.Fatalln(fmt.Errorf("error fetching project group: %s", name))
	}

	for _, p := range projectGroups {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, fmt.Errorf("project group '%s' could not be found, please check that it matches the project name exactly", name)
}

func GetProjectByName(client *od.Client, project string) (*projects.Project, error) {
	ps, err := client.Projects.Get(projects.ProjectsQuery{Name: project})
	if err != nil {
		return nil, fmt.Errorf("an error occured while fetching project '%s': %w", project, err)
	}

	for _, p := range ps.Items {
		if p.Name == project {
			return p, nil
		}
	}

	return nil, fmt.Errorf("project '%s' could not be found, please check that it matches the project name exactly", project)
}

func OctopusClient() *od.Client {
	apiURL, _ := url.Parse(viper.GetString("url"))

	client, err := od.NewClient(nil, apiURL, viper.GetString("apikey"), "")
	if err != nil {
		log.Fatalf("error creating Octopus Api client: %v", err)
	}

	return client
}
