package github

import (
	"fmt"
	"strings"

	"github.com/eunomie/github-board-exporter/configuration"
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
	Cards       []Card
	NbIssues    int
	Limit       int
	LimitSet    bool
	Exceeded    bool
	ExtraIssues int
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
	projectURLPattern         = "https://api.github.com/projects/%d"
	issuesMetricsPattern      = "board_issues{column=\"%s\",project=\"%d\"} %d"
	totalIssuesMetricsPattern = "board_issues_count{project=\"%d\"} %d"
	wipIssuesMetricsPattern   = "board_issues_wip{project=\"%d\"} %d"
	limitExceededPattern      = "board_limit_exceeded{column=\"%s\",project=\"%d\",exceeded=\"%t\",limit=\"%d\"} %d"
)

// NewProject creates a new representation of a github project
func NewProject(id int, github *Github, c *configuration.Configuration) (*Project, error) {
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
		col.NbIssues = col.numberOfIssues()
		col.Limit, col.LimitSet = c.Limit(col.Name)
		if col.LimitSet {
			col.Exceeded = col.NbIssues > col.Limit
			col.ExtraIssues = max(0, col.NbIssues-col.Limit)
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
	totalIssues := 0
	wipIssues := 0
	cols := len(p.Columns)
	for i, col := range p.Columns {
		nbIssues := col.numberOfIssues()
		totalIssues += nbIssues
		if i > 0 && i < cols-1 {
			wipIssues += nbIssues
		}
		metric := fmt.Sprintf(issuesMetricsPattern, col.Name, p.ID, nbIssues)
		metrics = append(metrics, metric)
		if col.LimitSet {
			limitMetric := fmt.Sprintf(limitExceededPattern, col.Name, p.ID,
				col.Exceeded, col.Limit, col.ExtraIssues)
			metrics = append(metrics, limitMetric)
		}
	}
	total := fmt.Sprintf(totalIssuesMetricsPattern, p.ID, totalIssues)
	metrics = append(metrics, total)
	wip := fmt.Sprintf(wipIssuesMetricsPattern, p.ID, wipIssues)
	metrics = append(metrics, wip)

	return strings.Join(metrics, "\n")
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
