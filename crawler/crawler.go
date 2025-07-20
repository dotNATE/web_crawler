package crawler

import (
	"fmt"
	"net/http"
	"net/url"
)

type Crawler struct {
	Base *url.URL

	HttpClient *http.Client
}

func NewCrawler(base string) (*Crawler, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		Base:       u,
		HttpClient: &http.Client{},
	}, nil
}

func (c *Crawler) Crawl() {
	fmt.Println("Crawl invoked")
}
