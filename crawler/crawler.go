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

	Results      map[string][]string
	ResultsMutex sync.Mutex
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
		Results:    make(map[string][]string),
	}, nil
}

func (c *Crawler) Start() {
	c.Crawl(c.Base)
	c.WaitGroup.Wait()
}

func (c *Crawler) Crawl(u *url.URL) {
	normalizedUrl := NormaliseURL(u).String()
	if _, loaded := c.Visited.LoadOrStore(normalizedUrl, true); loaded {
		return
	}
	fmt.Printf("Visiting: %s\n", normalizedUrl)

	c.WaitGroup.Add(1)

	go func() {
		defer c.WaitGroup.Done()

		c.Sem <- struct{}{}
		defer func() {
			<-c.Sem
		}()

		resp, err := c.HttpClient.Get(normalizedUrl)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		links, err := ExtractLinks(u, resp.Body)
		if err != nil {
			return
		}

		var linkStrings []string
		for _, link := range links {
			fmt.Printf("  -> %s\n", link.String())
			linkStrings = append(linkStrings, link.String())
		}

		c.ResultsMutex.Lock()
		c.Results[normalizedUrl] = linkStrings
		c.ResultsMutex.Unlock()

		for _, link := range links {
			if u.Host == c.Base.Host {
				c.Crawl(link)
			}
		}
	}()
}
