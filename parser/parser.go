package parser

import (
	"golang.org/x/net/html"
)


func Parse() {
	html.Parse(nil)
}

var walk func(*html.Node)


type Link struct {
	URL string
	Text string
}
