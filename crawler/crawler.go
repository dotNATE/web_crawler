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

func (c *Crawler) Start() {
	c.Crawl(c.Base)
}

func (c *Crawler) Crawl(u *url.URL) {
	fmt.Println("Crawl invoked")

	// normalise url

	// print to console

	// get content from page

	// extract links

	// print links to console

	// FOR LINKS
	// determine if link should be visited (e.g. is internal?)

	// recursively invoke crawl?
	// END FOR
}
