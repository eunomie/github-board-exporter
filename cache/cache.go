package cache

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Cache store a content during some caching duration
type Cache struct {
	caching time.Duration
	refresh func() (string, error)
	Content string
}

// NewCache creates a new cache object with refresh method
func NewCache(caching time.Duration, refresh func() (string, error)) *Cache {
	c := Cache{
		caching: caching,
		refresh: refresh,
		Content: "",
	}
	c.autoRefreshContent()
	go func() {
		c.refreshContent()
	}()
	return &c
}

func (c *Cache) autoRefreshContent() {
	ticker := time.NewTicker(c.caching)
	go func() {
		for _ = range ticker.C {
			if err := c.refreshContent(); err != nil {
				log.Errorf("could not auto refresh content: %v", err)
			}
		}
	}()
}

func (c *Cache) refreshContent() error {
	log.Println("refresh content...")

	content, err := c.refresh()
	if err != nil {
		return fmt.Errorf("could not refresh content: %v", err)
	}
	c.Content = content
	log.Println("...refreshed")

	return nil
}
