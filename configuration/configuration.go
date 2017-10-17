package configuration

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	nbDevMetricPattern = "github_board_number_of_developers %d"
)

// Configuration contains project id, access token to be authentified and user.
type Configuration struct {
	AccessToken string
	ProjectID   int    `yaml:"project_id"`
	Org         string `yaml:"github_org"`
	Repo        string `yaml:"github_repo"`
	NbDevs      int    `yaml:"number_of_developers"`
	Columns     []Column
	Limits      map[string]Limit
	Wip         map[string]bool
}

// Column contains configuration for a column as limit
type Column struct {
	Name  string
	Limit int
	Wip   bool
}

// Limit contains limit configuration and set attribute for a column
type Limit struct {
	Limit int
	Set   bool
}

// NewConfiguration reads config,yaml
func NewConfiguration() (*Configuration, error) {
	accessToken, set := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !set {
		return nil, fmt.Errorf("GITHUB_ACCESS_TOKEN must be defined")
	}

	conf := Configuration{AccessToken: accessToken, Limits: map[string]Limit{}, Wip: map[string]bool{}}

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

	if conf.Org == "" {
		return nil, fmt.Errorf("github user is missing")
	}

	if conf.Repo == "" {
		return nil, fmt.Errorf("github repo is missing")
	}

	if conf.NbDevs == 0 {
		return nil, fmt.Errorf("number_of_developers is missing")
	}

	for _, col := range conf.Columns {
		conf.Wip[col.Name] = col.Wip
		conf.Limits[col.Name] = Limit{col.Limit, col.Limit > 0}
	}

	return &conf, nil
}

// Metrics returns conf metrics
func (c *Configuration) Metrics() string {
	return fmt.Sprintf(nbDevMetricPattern, c.NbDevs)
}
