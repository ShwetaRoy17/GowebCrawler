package parser

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	URL  string
	Text string
}

func Parse(base *url.URL, body string) ([]Link, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	var links []Link

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string

			for _, a := range n.Attr {
				if a.Key == "href" {
					href = a.Val
					break
				}
			}
			if href != "" {
				if absURL := resolveURL(base, href); absURL != "" {

					links = append(links, Link{
						URL:  href,
						Text: strings.TrimSpace(nodeText(n)),
					})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
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
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(nodeText(c))
	}
	return b.String()
}

func resolveURL(base *url.URL, href string) string {
	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	resolved := base.ResolveReference(ref).String()
	if checkValidUrl(resolved) {
		return resolved
	}
	return ""
}

func checkValidUrl(url string) bool {
	if strings.HasPrefix(url, "mailto:") || strings.HasPrefix(url, "tel:") || strings.HasPrefix(url, "data:") || strings.HasPrefix(url, "ftp:") {
		return false
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}
