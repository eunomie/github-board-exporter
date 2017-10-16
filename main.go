package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eunomie/github-board-exporter/cache"
	"github.com/eunomie/github-board-exporter/configuration"
	"github.com/eunomie/github-board-exporter/github"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf, err := configuration.NewConfiguration()
	if err != nil {
		log.Fatalf("could not read configuration: %v", err)
	}

	log.Printf("project id %d", conf.ProjectID)

	g, err := github.NewGithub(conf.AccessToken)
	if err != nil {
		log.Fatalf("could not create Github client: %v", err)
	}

	c := cache.NewCache(30*time.Minute, allMetrics(conf, g))
	http.HandleFunc("/metrics", metrics(c))
	http.ListenAndServe(":8080", nil)
}

func metrics(c *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("ask metrics")
		content := c.Content
		fmt.Fprintln(w, content)
	}
}

func allMetrics(c *configuration.Configuration, g *github.Github) func() (string, error) {
	return func() (string, error) {
		id := c.ProjectID
		u := c.User

		p, err := github.NewProject(id, g, c)
		if err != nil {
			return "", fmt.Errorf("could not read project %d: %v", id, err)
		}
		boardMetrics := p.Metrics()

		pr, err := github.PullRequestsMetrics(g, u)
		if err != nil {
			return "", fmt.Errorf("could not read pull request metrics for user %s: %v", u, err)
		}

		metrics := fmt.Sprintf("%s\n%s", boardMetrics, pr)
		return metrics, nil
	}
}
