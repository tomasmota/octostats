package environment

import (
	"fmt"
	"log"
	"path/filepath"

	od "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
)

func GetEnvironmentIDs(client *od.Client, blob string) ([]string, error) {

	envs, err := client.Environments.GetAll()
	if err != nil {
		log.Fatalln(fmt.Errorf("error fetching the environment from octopus: %w", err))
	}

	if len(envs) == 0 {
		return nil, fmt.Errorf("environment matching '%s' not found", blob)
	}

	matches := make([]string, 0)
	for _, env := range envs {
		match, _ := filepath.Match(blob, env.Name)
		if match {
			matches = append(matches, env.ID)
		}
	}

	return matches, nil
}

func GetEnvironmentID(client *od.Client, environmentName string) string {
	ids, err := GetEnvironmentIDs(client, environmentName)
	if err != nil {
		log.Fatalln(fmt.Errorf("error fetching ID for environment: %s", environmentName))
	} else if len(ids) == 0 {
		log.Fatalln(fmt.Errorf("no environment found matching: %s", environmentName))
	} else if len(ids) > 1 {
		log.Fatalln(fmt.Errorf("found more than one environment matching: %s", environmentName))
	}
	return ids[0]
}
