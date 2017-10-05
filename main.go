package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eunomie/github-board-exporter/github"
	log "github.com/sirupsen/logrus"
)

const (
	issueMetricsPattern = "issues{column=\"%s\"} %d"
)

func main() {
	conf, err := newConfiguration()
	if err != nil {
		log.Fatalf("could not read configuration: %v", err)
	}

	log.Printf("Read project configuration for id %d", conf.projectID)

	g, err := github.NewGithub(conf.accessToken)
	if err != nil {
		log.Fatalf("could not create Github client: %v", err)
	}

	m := metrics(g, conf.projectID)
	http.HandleFunc("/metrics", m)
	http.ListenAndServe(":8080", nil)
}

func metrics(g *github.Github, id int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := github.NewProject(id, g)
		if err != nil {
			errString := fmt.Sprintf("could not read project %d: %v", id, err)
			log.Errorln(errString)
			http.Error(w, errString, http.StatusInternalServerError)
		}
		metrics := []string{}
		for _, col := range p.Columns {
			metric := fmt.Sprintf(issueMetricsPattern, col.Name, col.NumberOfIssues())
			metrics = append(metrics, metric)
		}
		fmt.Fprintln(w, strings.Join(metrics, "\n"))
	}
}
