package configuration

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	nbDevMetricPattern = "number_of_developers %d"
)

// Configuration contains project id, access token to be authentified and user.
type Configuration struct {
	AccessToken string
	ProjectID   int    `yaml:"project_id"`
	User        string `yaml:"github_user"`
	NbDevs      int    `yaml:"number_of_developers"`
	Limits      []Limit
}

// Limit is the maximum number of task per column
type Limit struct {
	Name  string
	Limit int
}

// NewConfiguration reads config,yaml
func NewConfiguration() (*Configuration, error) {
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

	if conf.NbDevs == 0 {
		return nil, fmt.Errorf("number_of_developers is missing")
	}

	return &conf, nil
}

// Limit returns the maximum number of items in a column
func (c *Configuration) Limit(colName string) (int, bool) {
	for _, limit := range c.Limits {
		if limit.Name == colName {
			return limit.Limit, true
		}
	}
	return 0, false
}

// Metrics returns conf metrics
func (c *Configuration) Metrics() string {
	return fmt.Sprintf(nbDevMetricPattern, c.NbDevs)
}
