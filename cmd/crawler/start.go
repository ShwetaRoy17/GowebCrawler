package main

import (
	"fmt"
	"net/url"
	
	"github.com/ShwetaRoy17/GowebCrawler/internal/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/internal/parser"
	"github.com/spf13/cobra"
)

var (
	seedUrl string
	depth int
	concurrency int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the web crawler",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run(seedUrl, depth, concurrency); err != nil {
			return fmt.Errorf("failed to start crawler: %w", err)
		}
		return nil
	},
}

func run(seed string, depth int, concurrency int) error {
	fmt.Printf("Crawling %s , depth: %d, concurrency: %d\n",seed,depth, concurrency)
	body, err := fetcher.Fetch(seed)
	if err != nil {
		return fmt.Errorf("error fetching data %s:%w",seed, err)

	}
	parsedUrl, err := url.Parse(seed)
	if err != nil {
		return fmt.Errorf("error parsing url %s:%w",seed, err)
	}
	links, err := parser.Parse(parsedUrl, body)
	if err != nil {
		return fmt.Errorf("error parsing body %s:%w",seed, err)
	}
	fmt.Printf("Found %d links on %s\n", len(links), seed)
	return nil
}


func init() {
	startCmd.Flags().StringVar(&seedUrl, "seed", "", "seed URL to crawl")
	startCmd.Flags().IntVar(&depth, "depth",3, "depth of links to be followd")
	startCmd.Flags().IntVar(&concurrency, "concurrency",10,"no. of goroutines")
	startCmd.MarkFlagRequired("seed")
}