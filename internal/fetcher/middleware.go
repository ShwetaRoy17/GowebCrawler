package fetcher

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type LoggingMiddleware struct {
	next Doer
}

type Metrics struct {
	TotalRequests   atomic.Int64
	SuccessRequests atomic.Int64
	FailedRequests  atomic.Int64
	TotalBytes      atomic.Int64
}

type MetricsMiddleware struct {
	next    Doer
	metrics *Metrics
}

type CompressionMiddleware struct {
	next Doer
}

func NewLoggingMiddleware(next Doer) *LoggingMiddleware {
	return &LoggingMiddleware{next: next}
}

func NewMetricsMiddleware(next Doer) *MetricsMiddleware {
	return &MetricsMiddleware{
		next:    next,
		metrics: &Metrics{},
	}
}

func NewCompressionMiddleware(next Doer) *CompressionMiddleware {
	return &CompressionMiddleware{next: next}
}

func (m *MetricsMiddleware) Do(req *http.Request) (*http.Response, error) {
	m.metrics.TotalRequests.Add(1)
	res, err := m.next.Do(req)
	if err != nil {
		m.metrics.FailedRequests.Add(1)
		return nil, err
	}
	if res.StatusCode == http.StatusOK {
		m.metrics.SuccessRequests.Add(1)
		m.metrics.TotalBytes.Add(res.ContentLength)
	}
	return res, nil
}

func (m *LoggingMiddleware) Do(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	log.Printf("-> %s %s", req.Method, req.URL)
	res, err := m.next.Do(req)

	if err != nil {
		log.Printf("x %s %s error: %v", req.Method, req.URL, err)
		return nil, err
	}
	log.Printf("<- %d %s %s (%s) %s %s %s", res.StatusCode, req.Method, req.URL, time.Since(startTime), req.Method, req.URL, time.Since(startTime))
	return res, nil

}

func (c *CompressionMiddleware) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept-Encoding", "gzip")
	res, err := c.next.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, fmt.Errorf("creating gzip reader: %w", err)
		}
		res.Body = io.NopCloser(reader)
	}
	return res, nil
}
