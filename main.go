package main

import (
	"fmt"
	"github.com/ShwetaRoy17/GowebCrawler/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/parser"
)

func main() {
	body, err := fetcher.Fetch("https://www.google.com")
	if err != nil {
		fmt.Println("Error fetching the data:", err.Error())
		return
	}
	runes := []rune(body)
	if len(runes) > 500 {
		runes = runes[:500]
	}
	fmt.Println(string(runes))
	parser.Parse()
}
