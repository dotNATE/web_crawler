package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type Crawler struct {
	Base *url.URL

	HttpClient HTTPClient

	WaitGroup sync.WaitGroup
	Visited   sync.Map
	Sem       chan struct{}
}

type HTTPClient interface {
	Get(string) (*http.Response, error)
}

func NewCrawler(base string, httpClient HTTPClient, concurrencyLimit int) (*Crawler, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		Base:       u,
		HttpClient: httpClient,
		Sem:        make(chan struct{}, concurrencyLimit),
	}, nil
}

func (c *Crawler) Start() {
	c.Crawl(c.Base)
	c.WaitGroup.Wait()
}

func (c *Crawler) Crawl(u *url.URL) {
	normalizedUrl := NormaliseURL(u)
	if _, loaded := c.Visited.LoadOrStore(normalizedUrl.String(), true); loaded {
		return
	}
	fmt.Printf("Visiting: %s\n", normalizedUrl.String())

	c.WaitGroup.Add(1)

	go func() {
		defer c.WaitGroup.Done()

		c.Sem <- struct{}{}
		defer func() {
			<-c.Sem
		}()

		resp, err := c.HttpClient.Get(normalizedUrl.String())
		if err != nil {
			return
		}
		defer resp.Body.Close()

		links, err := ExtractLinks(u, resp.Body)
		if err != nil {
			return
		}
		fmt.Printf("LINK: %+v\n", links)

		for _, link := range links {
			fmt.Printf("  -> %s\n", link.String())

			if u.Host == c.Base.Host {
				c.Crawl(link)
			}
		}
	}()
}
