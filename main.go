package main

import (
	"github.com/ShwetaRoy17/GowebCrawler/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/parser"
)

func main() {
	body, err := fetcher.Get("https://www.google.com")
	if err != nil {
		println("Error fetching the data:", err.Error())
		return
	}
	println("Fetched data:", body)
	parser.Parse()
}
