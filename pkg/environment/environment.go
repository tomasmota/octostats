package environment

import (
	"fmt"
	"log"
	"path/filepath"

	od "github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
