package github

import (
	"fmt"
)

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
	URL      string `json:"url"`
	CardsURL string `json:"cards_url"`
	Name     string `json:"name"`
	Cards    *Cards
}

// Columns of a project
type Columns struct {
	github     *Github
	columnsURL string
	columns    []Column
}

// Card represents a card in a board column
type Card struct {
	URL        string `json:"url"`
	ID         int    `json:"id"`
	Note       string `json:"note"`
	ContentURL string `json:"content_url"`
}

// Cards of a column
type Cards struct {
	github   *Github
	cardsURL string
	cards    []Card
}

// Project from github
type Project struct {
	github *Github
	projectFields
	Columns *Columns
}

const (
	projectURLPattern = "https://api.github.com/projects/%d"
)

// NewProject creates a new representation of a github project
func NewProject(id int, github *Github) (*Project, error) {
	url := fmt.Sprintf(projectURLPattern, id)
	p := Project{github: github}
	if err := github.GetJSON(url, &p); err != nil {
		return nil, fmt.Errorf("could not get project for ID %d: %v", id, err)
	}
	cols, err := newColumns(&p, p.ColumnsURL)
	if err != nil {
		return nil, err
	}
	p.Columns = cols
	return &p, nil
}

func newColumns(p *Project, u string) (*Columns, error) {
	columns := Columns{github: p.github, columnsURL: u}
	if err := columns.fetch(); err != nil {
		return nil, err
	}
	return &columns, nil
}

// Columns fetch columns information from github
func (c *Columns) fetch() error {
	if err := c.github.GetJSON(c.columnsURL, &c.columns); err != nil {
		return fmt.Errorf("could not get columns from URL %s: %v", c.columnsURL, err)
	}

	for i := range c.columns {
		col := &c.columns[i]
		cards, err := newCards(c, col.CardsURL)
		if err != nil {
			return err
		}
		col.Cards = cards
	}
	return nil
}

// Count the number of columns
func (c *Columns) Count() int {
	if err := c.fetch(); err != nil {
		return 0
	}
	return len(c.columns)
}

// Get a column by index
func (c *Columns) Get(idx int) (*Column, error) {
	if idx < 0 || idx >= len(c.columns) {
		return nil, fmt.Errorf("out of range column for index %d", idx)
	}
	return &c.columns[idx], nil
}

// GetByName returns a column by his name
func (c *Columns) GetByName(name string) (*Column, error) {
	for i := range c.columns {
		col := &c.columns[i]
		if col.Name == name {
			return col, nil
		}
	}
	return nil, fmt.Errorf("could not find column with name %s", name)
}

func newCards(c *Columns, u string) (*Cards, error) {
	cards := Cards{github: c.github, cardsURL: u}
	if err := cards.fetch(); err != nil {
		return nil, err
	}
	return &cards, nil
}

func (cards *Cards) fetch() error {
	if err := cards.github.GetJSON(cards.cardsURL, &cards.cards); err != nil {
		return fmt.Errorf("could not get cards from URL %s: %v", cards.cardsURL, err)
	}
	return nil
}

// Count the number of cards
func (cards *Cards) Count() int {
	return len(cards.cards)
}

// Get a card
func (cards *Cards) Get(idx int) (*Card, error) {
	if idx < 0 || idx >= len(cards.cards) {
		return nil, fmt.Errorf("out of range column for index %d", idx)
	}
	return &cards.cards[idx], nil
}
