package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration contains project id, access token to be authentified and user.
type Configuration struct {
	AccessToken string
	ProjectID   int    `yaml:"project_id"`
	User        string `yaml:"github_user"`
}

func newConfiguration() (*Configuration, error) {
	accessToken, set := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !set {
		return nil, fmt.Errorf("GITHUB_ACCESS_TOKEN must be defined")
	}

	conf := Configuration{AccessToken: accessToken}

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("could not parse configuration file: %v", err)
	}

	if conf.ProjectID == 0 {
		return nil, fmt.Errorf("project id is missing")
	}

	if conf.User == "" {
		return nil, fmt.Errorf("github user is missing")
	}

	return &conf, nil
}
