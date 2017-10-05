package main

import (
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

	project, err := github.NewProject(conf.projectID, g)
	if err != nil {
		log.Fatalf("could not read info for project %d: %v", conf.projectID, err)
	}

	log.Println(project.ID)

	log.Println(project.Columns.GetByName("📚 Backlog"))
	for i := 0; i < project.Columns.Count(); i++ {
		col, _ := project.Columns.Get(i)
		log.Println(col.Name)
		for j := 0; j < col.Cards.Count(); j++ {
			card, _ := col.Cards.Get(j)
			if card.Note != "" {
				log.Println("  " + card.Note)
			} else {
				log.Println("  " + card.ContentURL)
			}
		}
	}
}
