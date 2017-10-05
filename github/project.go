package github

import (
	"fmt"
	"strings"
)

// Project from github
type Project struct {
	projectFields
	Columns []Column
}

type projectFields struct {
	URL        string `json:"url"`
	HTMLURL    string `json:"html_url"`
	ColumnsURL string `json:"columns_url"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Body       string `json:"body"`
}

// Column represents a project column from the Github API
type Column struct {
	columnFields
	Cards []Card
}

type columnFields struct {
	URL      string `json:"url"`
	CardsURL string `json:"cards_url"`
	Name     string `json:"name"`
}

// Card represents a card in a board column
type Card struct {
	cardFields
}

type cardFields struct {
	URL        string `json:"url"`
	ID         int    `json:"id"`
	Note       string `json:"note"`
	ContentURL string `json:"content_url"`
}

const (
	projectURLPattern   = "https://api.github.com/projects/%d"
	issueMetricsPattern = "issues{column=\"%s\"} %d"
)

// NewProject creates a new representation of a github project
func NewProject(id int, github *Github) (*Project, error) {
	url := fmt.Sprintf(projectURLPattern, id)
	p := Project{}

	if err := github.GetJSON(url, &p); err != nil {
		return nil, fmt.Errorf("could not fetch project %d: %v", id, err)
	}

	if err := github.GetJSON(p.ColumnsURL, &p.Columns); err != nil {
		return nil, fmt.Errorf("could not fetch columns for project %d: %v", id, err)
	}

	for i := range p.Columns {
		col := &p.Columns[i]
		if err := github.GetJSON(col.CardsURL, &col.Cards); err != nil {
			return nil, fmt.Errorf("could not fetch cards for project %d: %v", id, err)
		}
	}

	return &p, nil
}

// NumberOfIssues count the number of cards with a content (issues, PR)
// in a column
func (c *Column) numberOfIssues() int {
	n := 0
	for _, card := range c.Cards {
		if card.ContentURL != "" {
			n++
		}
	}
	return n
}

// Metrics compatible with prometheus
func (p *Project) Metrics() string {
	metrics := []string{}
	for _, col := range p.Columns {
		metric := fmt.Sprintf(issueMetricsPattern, col.Name, col.numberOfIssues())
		metrics = append(metrics, metric)
	}
	return strings.Join(metrics, "\n")
}
