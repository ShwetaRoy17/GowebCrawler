package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"crypto/tls"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

type Options struct {
    Timeout   time.Duration
    UserAgent string
    SkipTLS   bool
}

type Fetcher struct {
    client    *http.Client
    userAgent string
}

func NewFetcher(options Options) *Fetcher {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.SkipTLS,
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   options.Timeout,
	}
	return &Fetcher{
		client:    client,
		userAgent: options.UserAgent,
	}
}


func(f *Fetcher) Fetch(url string) (string, error) {
	req, err  := http.NewRequest(http.MethodGet,url,nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent",f.userAgent)

	res, err := f.client.Do(req)

	if err != nil {
		return "",fmt.Errorf("network error: %w",err)
	}

	defer func(){
		if err := res.Body.Close(); err!= nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil

}
