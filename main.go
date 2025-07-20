package main

import (
	"fmt"
	"net/http"

	"example.com/crawler"
)

func main() {
	c, err := crawler.NewCrawler("https://monzo.com", &http.Client{}, 5)
	if err != nil {
		fmt.Printf("Failed to initalise Crawler: %+v", err)
		return
	}

	c.Start()
	fmt.Printf("\n>>>> Pages Visited %d\n", len(c.Results))
}
