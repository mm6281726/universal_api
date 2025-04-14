package models

import (
	"time"
)

// APIDocRequest represents a request to scrape an API documentation
type APIDocRequest struct {
	URL         string `json:"url" binding:"required"`
	Description string `json:"description"`
}

// APIDoc represents a scraped API documentation
type APIDoc struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Endpoints   []Endpoint `json:"endpoints"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Endpoint represents an API endpoint
type Endpoint struct {
	Path        string      `json:"path"`
	Method      string      `json:"method"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters"`
	Responses   []Response  `json:"responses"`
}

// Parameter represents an API endpoint parameter
type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"` // query, path, header, body
	Required    bool   `json:"required"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Response represents an API endpoint response
type Response struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	Schema      string `json:"schema,omitempty"` // JSON schema as string
}


