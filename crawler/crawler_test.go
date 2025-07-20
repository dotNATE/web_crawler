package crawler_test

import (
	"testing"

	"example.com/crawler"
)

func TestCrawl(t *testing.T) {
	c, err := crawler.NewCrawler("https://example.com")
	if err != nil {
		t.Errorf("Failed to initialise Crawler: error: %+v", err)
	}

	c.Crawl()
}
