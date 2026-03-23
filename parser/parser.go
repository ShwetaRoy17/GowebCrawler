package parser

import (
	"strings"

	"golang.org/x/net/html"
)


func Parse() {
	html.Parse(nil)
}




type Link struct {
	URL string
	Text string
}

func extractLinks(body string)([]Link, error ){
	doc, err := html.Parse(strings.NewReader(body))
	if err!= nil {
		return nil, err
	}
	var links []Link

	var walk func(*html.Node)
	walk = func(n *html.Node){
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string

			for _, a := range n.Attr {
				if a.Key == "href"{
					href = a.Val
					break
				}
			}
			if href != "" {
				links = append(links, Link{
					URL: href,
					Text: strings.TrimSpace(nodeText(n)),
				})
			}
		}
		for c:= n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return links, nil
}

func nodeText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var b strings.Builder
	for c:= n.FirstChild; c!= nil ; c = c.NextSibling {
		b.WriteString(nodeText(c))
	}
	return b.String()
}
