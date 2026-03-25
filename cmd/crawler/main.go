package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/ShwetaRoy17/GowebCrawler/internal/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/internal/parser"
)

func h1(urlString string){
		fmt.Println(urlString,time.Now())

	body, err := fetcher.Fetch(urlString)
	if err != nil {
		fmt.Println("Error fetching the data:", err.Error())
		return
	}
	parsedUrl,err := url.Parse(urlString)
	
	abc,err := parser.Parse(parsedUrl,body)
	if err != nil {
		fmt.Println("Error parsing the data:", err.Error())
		return
	}
	fmt.Printf("Found %d links\n", abc)
	fmt.Println(urlString,time.Now())

}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.twitter.com",
		"https://www.linkedin.com",
		"https://www.github.com",
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			h1(u)
		}(url)
	}
	wg.Wait()

	// runes := []rune(body)
	// if len(runes) > 500 {
	// 	runes = runes[:500]
	// }
	// fmt.Println(string(runes))
}
