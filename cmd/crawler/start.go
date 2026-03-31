package main

import (
	"fmt"
	"net/url"
	"time"
	"sync"

	"github.com/spf13/cobra"

	"github.com/ShwetaRoy17/GowebCrawler/internal/config"
	"github.com/ShwetaRoy17/GowebCrawler/internal/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/internal/parser"
)

var (
	seedUrl     string
	depth       int
	concurrency int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the web crawler",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration loading failed :%w", err)
		}
		fetcher := fetcher.NewFetcher(fetcher.Options{
			Timeout:   10 * time.Second,
			UserAgent: cfg.UserAgent,
			SkipTLS:   true,
		})
		cfg.SeedUrl = seedUrl
		return run(*cfg,fetcher)
	},
}

func run(cfg config.Config, fetcher *fetcher.Fetcher) error {
	visited := make(map[string]bool)
	var mu sync.Mutex
	var crawlerF func(urlp string, currd int) error
	crawlerF = func(urlp string, currd int) error {
		if currd > cfg.MaxDepth {
			return nil
		}
		mu.Lock()
		if visited[urlp] {
			mu.Unlock()
			return nil
		}
		visited[urlp] = true
		mu.Unlock()

		body, err := fetcher.FetchWithRetry(urlp,3)
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", urlp, err)
		}

		parsedUrl, err := url.Parse(urlp)
		if err != nil {
			return fmt.Errorf("failed to parse URL %s: %w", urlp, err)
		}
		links, err := parser.Parse(parsedUrl, body)
		if err != nil {
			return fmt.Errorf("failed to parse links %w", err)
		}
		fmt.Printf("crawled: %s, depth: %d, concurrency:%d\n", urlp, cfg.MaxDepth, cfg.Concurrency)
		var wg sync.WaitGroup
		sem := make(chan struct{}, cfg.Concurrency)
		for _, link := range links {
			wg.Add(1)
			sem <- struct{}{}
			go func(li parser.Link) {
				defer wg.Done()
				defer func() { <-sem }()
				if err := crawlerF(li.URL, currd+1); err != nil {
					fmt.Printf("error crawling %s: %v\n", li.URL, err)
				}

			}(link)
		}
		wg.Wait()
		return nil
	}
	return crawlerF(cfg.SeedUrl, 0)
}

func init() {
	startCmd.Flags().StringVar(&seedUrl, "seed", "", "seed URL to crawl")
	startCmd.Flags().IntVar(&depth, "depth", 3, "depth of links to be followd")
	startCmd.Flags().IntVar(&concurrency, "concurrency", 10, "no. of goroutines")
	if err := startCmd.MarkFlagRequired("seed"); err != nil {
		panic(err)
	}

}
