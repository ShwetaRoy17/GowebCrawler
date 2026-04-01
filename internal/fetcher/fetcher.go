package fetcher

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Options struct {
	Timeout   time.Duration
	UserAgent string
	SkipTLS   bool
	RateLimit rate.Limit
	Burst     int
}

type Fetcher struct {
	client    *http.Client
	userAgent string
	limiter   *rate.Limiter
}

type HTTPError struct {
	StatusCode int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http error: %d", e.StatusCode)
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
		limiter:   rate.NewLimiter(options.RateLimit, options.Burst),
	}
}

func (f *Fetcher) Fetch(url string) (string, error) {
	if err := f.limiter.Wait(context.Background()); err != nil {
		return "", fmt.Errorf("rate limit error: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", f.userAgent)

	res, err := f.client.Do(req)

	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return "", &HTTPError{StatusCode: res.StatusCode}
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil

}

func (f *Fetcher) FetchWithRetry(url string, retries int) (string, error) {
	var lastErr error
	for attempt := 0; attempt < retries; attempt++ {
		var body string
		body, lastErr := f.Fetch(url)
		if lastErr == nil {
			return body, nil
		}
		var httpErr HTTPError
		if errors.As(lastErr, &httpErr) {
			return "", lastErr
		}
		wait := time.Duration(1<<attempt) * time.Second
		time.Sleep(wait)
	}
	return "", fmt.Errorf("all %d attempts failed: %w", retries, lastErr)

}
