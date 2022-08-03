package util

import (
	"fmt"
	"log"
	"net/url"

	od "github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/spf13/viper"
)

func GetProjectByName(client *od.Client, project string) (*od.Project, error) {
	ps, err := client.Projects.Get(od.ProjectsQuery{Name: project})
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
