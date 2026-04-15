package main


import (
	"fmt"
	"sync"

	"net/http"
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

func main(){
	// server := NewServer()

	http.HandleFunc("/hello",func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	if err:= http.ListenAndServe(":8080",nil); err != nil {
		fmt.Printf("Server error : %v\n",err)
	}
}