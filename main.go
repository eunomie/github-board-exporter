package main

import (
	"fmt"

	"github.com/eunomie/github-board-exporter/github"
	log "github.com/sirupsen/logrus"
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

	project, err := g.GetString(fmt.Sprintf("https://api.github.com/projects/%d", conf.projectID))
	if err != nil {
		log.Fatalf("could not read info for project %d: %v", conf.projectID, err)
	}
	log.Println(project)
}
