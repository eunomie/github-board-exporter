package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eunomie/github-board-exporter/cache"
	"github.com/eunomie/github-board-exporter/github"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf, err := newConfiguration()
	if err != nil {
		log.Fatalf("could not read configuration: %v", err)
	}

	log.Printf("project id %d", conf.projectID)

	g, err := github.NewGithub(conf.accessToken)
	if err != nil {
		log.Fatalf("could not create Github client: %v", err)
	}

	c := cache.NewCache(30*time.Minute, allMetrics(conf.projectID, conf.user, g))
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

func allMetrics(id int, u string, g *github.Github) func() (string, error) {
	return func() (string, error) {
		p, err := github.NewProject(id, g)
		if err != nil {
			return "", fmt.Errorf("could not read project %d: %v", id, err)
		}

		pr, err := github.PullRequestsMetrics(g, u)
		if err != nil {
			return "", fmt.Errorf("could not read pull request metrics for user %s: %v", u, err)
		}

		metrics := fmt.Sprintf("%s\n%s", p.Metrics(), pr)
		return metrics, nil
	}
}
