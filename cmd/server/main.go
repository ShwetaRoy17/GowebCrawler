package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"os"
	"context"

	"github.com/ShwetaRoy17/GowebCrawler/internal/database"
	"github.com/ShwetaRoy17/GowebCrawler/internal/config"
	"github.com/ShwetaRoy17/GowebCrawler/internal/fetcher"
	"github.com/ShwetaRoy17/GowebCrawler/internal/parser"
	"github.com/ShwetaRoy17/GowebCrawler/internal/models"
	"github.com/joho/godotenv"
)

type StartCrawlRequest struct {
	Seed        string `json:"seed"`
	Depth       int    `json:"depth"`
	Concurrency int    `json:"concurrency"`
}



type Server struct {
	mu   sync.Mutex
	db *database.DB
}

func NewServer(db *database.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) handleStartCrawl(w http.ResponseWriter, r *http.Request) {

	var req StartCrawlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	if req.Seed == "" {
		http.Error(w, "missing seed parameter", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	seed := req.Seed
	if seed == "" {
		http.Error(w, "missing seed parameter", http.StatusBadRequest)
		return
	}

	n,err := s.db.LengthJobs(context.Background())
	if err != nil {
		http.Error(w, "failed to count jobs", http.StatusInternalServerError)
	
		return
	}

	jobID := fmt.Sprintf("job-%d", n+1)
	job := &models.JobStatus{
		ID:          jobID,
		Status:      "in_progress",
		Seed:        seed,
		Depth:       req.Depth,
		Concurrency: req.Concurrency,
	}

	if err:= s.db.InsertJob(context.Background(),job); err!= nil {
		http.Error(w, "failed to save job to database", http.StatusInternalServerError)
		return
	}
	go s.runCrawl(job, seed)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (s *Server) runCrawl(job *models.JobStatus, seed string) {
	cfg, err := config.Load()
	if err != nil {
		s.mu.Lock()
		job.Status = "Failed"
		job.Error = err.Error()
		s.mu.Unlock()
		if err := s.db.UpdateJob(context.Background(),job);err != nil {
			fmt.Printf("Error updating job in database: %v\n", err)
		}

		return
	}
	f := fetcher.NewFetcher(fetcher.Options{
		Timeout:   cfg.Timeout,
		UserAgent: cfg.UserAgent,
		RateLimit: cfg.RateLimit,
		Burst:     cfg.Burst,
	})

	visited := make(map[string]bool)
	var mu sync.Mutex

	var crawl func(url string, depth int)
	crawl = func(urlS string, depth int) {
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

		body, err := f.FetchWithRetry(urlS, 3)
		if err != nil {
			return
		}
		parsedUrl, err := url.Parse(urlS)
		if err != nil {
			return
		}

		links, err := parser.Parse(parsedUrl, body)
		if err != nil {
			return
		}

		mu.Lock()
		job.Pages++
		mu.Unlock()
		var wg sync.WaitGroup
		for _, link := range links {
			wg.Add(1)
			go crawl(link.URL, depth+1)
		}
		wg.Wait()
		defer wg.Done()
	}

	crawl(seed, 0)

	s.mu.Lock()
	job.Status = "completed"
	s.mu.Unlock()
	if err := s.db.UpdateJob(context.Background(),job);err != nil {
		fmt.Printf("Error updating job in database: %v\n", err)
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	job, err := s.db.GetJob(context.Background(), jobID)
	if err != nil {
		http.Error(w, "error fetching job from database", http.StatusInternalServerError)
		return
	}
	if job == nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request) {
	
	DBjobs,err := s.db.ListJobs(context.Background())
	
	if err != nil {
		http.Error(w, "error fetching jobs from database", http.StatusInternalServerError)
		return
	}
	jobs := make([]*models.JobStatus, len(DBjobs))
	for i, job := range DBjobs {
		jobs[i] = job
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found)")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL not set in environment")
		return
	}
	db, err := database.New(dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	if err := db.CreateSchema(context.Background()); err != nil {
		fmt.Printf("Error creating database schema: %v\n", err)
		os.Exit(1)
	}
	server := NewServer(db)

	http.HandleFunc("/start", server.handleStartCrawl)
	http.HandleFunc("/status", server.handleStatus)
	http.HandleFunc("/jobs",server.handleListJobs)

	fmt.Println("server starting on : 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error : %v\n", err)
	}
}
