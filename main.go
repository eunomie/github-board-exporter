package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eunomie/github-board-exporter/github"
	log "github.com/sirupsen/logrus"
)

type cache struct {
	RefreshedAt time.Time
	Metrics     string
}

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

	c := cache{}

	m := metrics(g, conf.projectID, &c)
	http.HandleFunc("/metrics", m)
	http.ListenAndServe(":8080", nil)
}

func metrics(g *github.Github, id int, c *cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("ask metrics")
		if c.Metrics == "" || c.RefreshedAt.Before(time.Now().Add(-30*time.Minute)) {
			log.Println("  fetch project")
			c.RefreshedAt = time.Now()

			p, err := github.NewProject(id, g)
			if err != nil {
				errString := fmt.Sprintf("could not read project %d: %v", id, err)
				log.Errorln(errString)
				http.Error(w, errString, http.StatusInternalServerError)
			}
			log.Println("  project fetched")
			c.Metrics = p.Metrics()
		}
		fmt.Fprintln(w, c.Metrics)
		log.Println("end metrics")
	}
}
