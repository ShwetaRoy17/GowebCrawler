package main

import (
	"github.com/ShwetaRoy17/GowebCrawler/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/parser"
)

func main() {
	fetcher.Fetch()
	parser.Parse()
}