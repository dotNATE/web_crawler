package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"example.com/crawler"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Not enough args provided to crawler. Expecting: url recursion_depth output_filename")
	}

	var maxDepth int
	if len(os.Args) > 2 && os.Args[2] != "" {
		var err error
		maxDepth, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Error parsing max concurrency depth input")
		}
	}

	c, err := crawler.NewCrawler(os.Args[1], &http.Client{}, 5, maxDepth)
	if err != nil {
		log.Fatalf("Failed to initalise Crawler: %+v", err)
	}

	c.Start()

	if len(os.Args) > 3 && os.Args[3] != "" {
		if err := c.ExportResults(os.Args[3]); err != nil {
			log.Fatalf("Failed to write results: %v", err)
		}
	}

	fmt.Printf("\n>>>> Pages Visited %d\n", len(c.Results))
}
