package main

import (
	"fmt"
	"os"
	"strconv"
)

// Configuration contains project id, access token to be authentified and user.
type Configuration struct {
	projectID   int
	accessToken string
	user        string
}

func newConfiguration() (*Configuration, error) {
	accessToken, set := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !set {
		return nil, fmt.Errorf("GITHUB_ACCESS_TOKEN must be defined")
	}
	projectIDStr, set := os.LookupEnv("PROJECT_ID")
	if !set {
		return nil, fmt.Errorf("PROJECT_ID must be defined")
	}
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse PROJECT_ID %s", projectIDStr)
	}
	user, set := os.LookupEnv("GITHUB_USER")
	if !set {
		return nil, fmt.Errorf("GITHUB_USER must be defined")
	}

	return &Configuration{
		projectID,
		accessToken,
		user,
	}, nil
}
