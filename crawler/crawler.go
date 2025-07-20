package crawler

import (
	"fmt"
	"net/http"
	"net/url"
)

type Crawler struct {
	Base *url.URL

	HttpClient HTTPClient
}

type HTTPClient interface {
	Get(string) (*http.Response, error)
}

func NewCrawler(base string, httpClient HTTPClient) (*Crawler, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		Base:       u,
		HttpClient: httpClient,
	}, nil
}

func (c *Crawler) Start() {
	c.Crawl(c.Base)
}

func (c *Crawler) Crawl(u *url.URL) {
	fmt.Println("Crawl invoked")

	// normalise url
	normalizedUrl := NormaliseURL(u)

	// print to console
	fmt.Println("Visiting:", u.String())

	// get content from page
	resp, err := c.HttpClient.Get(normalizedUrl.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// extract links
	links, err := ExtractLinks(u, resp.Body)
	if err != nil {
		return
	}

	// print links to console
	for _, link := range links {
		fmt.Printf("  -> %s\n", link.String())

		// determine if link should be visited (e.g. is internal?)
		if u.Host == c.Base.Host {
			// recursively invoke crawl?
			c.Crawl(link)
		}
	}
}
