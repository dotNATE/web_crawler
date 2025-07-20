package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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

	LogMutex sync.Mutex

	RecursionLimit int
}

type HTTPClient interface {
	Get(string) (*http.Response, error)
}

func NewCrawler(base string, httpClient HTTPClient, concurrencyLimit int, recursionLimit int) (*Crawler, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		Base:           u,
		HttpClient:     httpClient,
		Sem:            make(chan struct{}, concurrencyLimit),
		Results:        make(map[string][]string),
		RecursionLimit: recursionLimit,
	}, nil
}

func (c *Crawler) Start() {
	c.Crawl(c.Base, 0)
	c.WaitGroup.Wait()
}

func (c *Crawler) Crawl(u *url.URL, depth int) {
	if c.RecursionLimit > 0 && depth > c.RecursionLimit {
		return
	}

	normalizedUrl := NormaliseURL(u).String()
	if _, loaded := c.Visited.LoadOrStore(normalizedUrl, true); loaded {
		return
	}

	c.WaitGroup.Add(1)
	go func() {
		defer c.WaitGroup.Done()

		c.Sem <- struct{}{}
		defer func() {
			<-c.Sem
		}()

		resp, err := c.HttpClient.Get(u.String())
		if err != nil {
			c.LogMutex.Lock()
			fmt.Printf("Error fetching %s: %v\n", u.String(), err)
			c.LogMutex.Unlock()
			return
		}
		defer resp.Body.Close()

		links, err := ExtractLinks(u, resp.Body)
		if err != nil {
			c.LogMutex.Lock()
			fmt.Printf("Error extracting links from %s: %v\n", u.String(), err)
			c.LogMutex.Unlock()
			return
		}

		var linkStrings []string
		for _, link := range links {
			linkStrings = append(linkStrings, link.String())
		}

		c.ResultsMutex.Lock()
		c.Results[normalizedUrl] = linkStrings
		c.ResultsMutex.Unlock()

		c.LogMutex.Lock()
		fmt.Printf("Visited: %s\n", normalizedUrl)
		for _, link := range links {
			fmt.Printf("  -> %s\n", link.String())
		}
		c.LogMutex.Unlock()

		for _, link := range links {
			if c.ShouldVisit(NormaliseURL(link)) {
				c.Crawl(link, depth+1)
			}
		}
	}()
}

func (c *Crawler) ShouldVisit(u *url.URL) bool {
	return u.Host == c.Base.Host
}

func (c *Crawler) ExportResults(filename string) error {
	c.ResultsMutex.Lock()
	defer c.ResultsMutex.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c.Results)
}
