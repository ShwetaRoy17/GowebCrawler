package models

type JobStatus struct {
    ID          string `json:"id"`
    Status      string `json:"status"`
    Pages       int    `json:"pages"`
    Error       string `json:"error,omitempty"`
    Seed        string `json:"seed"`
    Depth       int    `json:"depth"`
    Concurrency int    `json:"concurrency"`
}