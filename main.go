package main

import (
	"fmt"
	"net/http"
	"strings"
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
	http.HandleFunc("/health", health(c))
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
		o := c.Org
		r := c.Repo
		metrics := []string{}

		p, err := github.NewProject(id, g, c)
		if err != nil {
			return "", fmt.Errorf("could not read project %d: %v", id, err)
		}
		metrics = append(metrics, p.Metrics(c))

		pr, err := github.IssuesMetrics(g, o, r)
		if err != nil {
			return "", fmt.Errorf("could not read pull request metrics for org %s and repo %s: %v", o, r, err)
		}
		metrics = append(metrics, pr)

		metrics = append(metrics, c.Metrics())

		return strings.Join(metrics, "\n"), nil
	}
}

func health(c *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if c.Content == "" {
			http.Error(w, "could not read cache", http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, "")
		}
	}
}
