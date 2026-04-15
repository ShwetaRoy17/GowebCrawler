package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/ShwetaRoy17/GowebCrawler/internal/config"
	"github.com/ShwetaRoy17/GowebCrawler/internal/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/internal/parser"
)

type JobStatus struct {
	ID string `json:"id"`
	Status string `json:"status"`
	Pages int `json:"pages"`
	Error string `json:"error,omitempty"`
}

type Server struct {
	jobs map[string]*JobStatus
	mu sync.Mutex
}

func NewServer() *Server {
	return &Server{
		jobs: make(map[string]*JobStatus),
	}
}

func (s *Server) handleStartCrawl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"method not allowed", http.StatusMethodNotAllowed)
	}

	seed := r.URL.Query().Get("seed")
	if seed == "" {
		http.Error(w , "missing seed parameter", http.StatusBadRequest)
		return
	}

	jobID := fmt.Sprintf("job-%d",len(s.jobs)+1)
	job := &JobStatus{
		ID: jobID,
		Status: "in_progress",
	}

	s.mu.Lock()
	s.jobs[jobID] = job
	s.mu.Unlock()

	go s.runCrawl(job, seed)

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(job)
}


func (s *Server) runCrawl(job *JobStatus, seed string) {
	cfg, err := config.Load()
	if err != nil {
		s.mu.Lock()
		job.Status = "Failed"
		job.Error = err.Error()
		s.mu.Unlock()
		return
	}
	 f := fetcher.NewFetcher(fetcher.Options{
		Timeout: cfg.Timeout,
		UserAgent: cfg.UserAgent,
		RateLimit: cfg.RateLimit,
		Burst: cfg.Burst,	
	 })

	 visited := make(map[string]bool)
	 var mu sync.Mutex

	 var crawl func(url string, depth int)
	 crawl = func(urlS string, depth int){
		if depth > cfg.MaxDepth {
			return
		}
		mu.Lock()
		if visited[urlS] {
			mu.Unlock()
			return
		}
		visited[urlS] = true
		mu.Unlock()

		body, err := f.FetchWithRetry(urlS,3)
		if err != nil {
			return
		}
		parsedUrl, err := url.Parse(urlS)
		if err != nil {
			return
		}

		links, err := parser.Parse(parsedUrl,body)
		if err != nil {
			return
		}
		
		mu.Lock()
		job.Pages++
		mu.Unlock()
		for _, link := range links {
			go crawl(link.URL, depth+1)
		}

	 }

	 crawl(seed,0)
	 
	 s.mu.Lock()
	 job.Status = "completed"
	 s.mu.Unlock()
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request){
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	job, ok := s.jobs[jobID]
	s.mu.Unlock()
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}


func main(){
	server := NewServer()

	http.HandleFunc("/start", server.handleStartCrawl)
	http.HandleFunc("/status", server.handleStatus)

	if err:= http.ListenAndServe(":8080",nil); err != nil {
		fmt.Printf("Server error : %v\n",err)
	}
}