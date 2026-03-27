package main

import (
	"fmt"
	"net/url"
	"sync"
	"github.com/spf13/cobra"
	"github.com/ShwetaRoy17/GowebCrawler/internal/config"
	"github.com/ShwetaRoy17/GowebCrawler/internal/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/internal/parser"
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
		cfg := config.Config{
			SeedUrl: seedUrl,
			MaxDepth: depth,
			Concurrency: concurrency,
		}
		return run(cfg)
	},
}

func run(cfg config.Config) error {
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

		body, err := fetcher.Fetch(urlp)
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", urlp, err)
		}

		parsedUrl, err := url.Parse(urlp)
		if err != nil {
			return fmt.Errorf("failed to parse URL %s: %w", urlp, err)
		}
		links,err := parser.Parse(parsedUrl, body)
		if err != nil {
			return fmt.Errorf("failed to parse links %s",&parsedUrl)
		}
		fmt.Printf("crawled: %s, depth: %d, found:%d\n", urlp,depth, len(links))
		var wg sync.WaitGroup
		sem := make(chan struct{}, cfg.Concurrency)
		for _, link := range links {
			wg.Add(1)
			sem <- struct{}{}
			go func(li parser.Link){
				defer wg.Done()
				defer func(){<-sem}()
				crawlerF(li.URL,currd+1)

			}(link)
		}
		wg.Wait()
		return nil
}
return crawlerF(cfg.SeedUrl, 0)
}


func init() {
	startCmd.Flags().StringVar(&seedUrl, "seed", "", "seed URL to crawl")
	startCmd.Flags().IntVar(&depth, "depth",3, "depth of links to be followd")
	startCmd.Flags().IntVar(&concurrency, "concurrency",10,"no. of goroutines")
	startCmd.MarkFlagRequired("seed")
}