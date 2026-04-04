package fetcher

import (
	"net/http"
	"time"
	"log"
)


type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type LoggingMiddleware struct {
	next Doer
}

func (m *LoggingMiddleware) Do(req *http.Request) (*http.Response, error){
	startTime := time.Now()
	log.Printf("-> %s %s",req.Method, req.URL)
	res, err := m.next.Do(req)

	if err != nil {
		log.Printf("x %s %s error: %v",req.Method, req.URL, err)
		return nil, err
	}
	log.Printf("<- %d %s %s (%s)", res.StatusCode, req.Method, req.URL, time.Since(startTime),req.Method, req.URL,time.Since(startTime))
	return res, nil

}

func NewLoggingMiddleware(next Doer) *LoggingMiddleware {
	return &LoggingMiddleware{next: next}
}

