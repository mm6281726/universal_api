package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIsSwaggerURL tests the isSwaggerURL function
func TestIsSwaggerURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com/swagger", true},
		{"https://example.com/api/swagger", true},
		{"https://example.com/openapi", true},
		{"https://example.com/api-docs", true},
		{"https://example.com/api/docs", false},
		{"https://example.com/documentation", false},
	}

	for _, test := range tests {
		result := isSwaggerURL(test.url)
		if result != test.expected {
			t.Errorf("isSwaggerURL(%s) = %v; expected %v", test.url, result, test.expected)
		}
	}
}

// TestIsRESTDocURL tests the isRESTDocURL function
func TestIsRESTDocURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com/api/doc", true},
		{"https://example.com/api/docs", true},
		{"https://example.com/api/reference", true},
		{"https://example.com/docs", false},
		{"https://example.com/reference", false},
		{"https://example.com/api", false},
	}

	for _, test := range tests {
		result := isRESTDocURL(test.url)
		if result != test.expected {
			t.Errorf("isRESTDocURL(%s) = %v; expected %v", test.url, result, test.expected)
		}
	}
}

// TestIsJSON tests the isJSON function
func TestIsJSON(t *testing.T) {
	tests := []struct {
		content  string
		expected bool
	}{
		{`{"name": "test"}`, true},
		{`[1, 2, 3]`, true},
		{`name: test`, false},
		{`<html></html>`, false},
	}

	for _, test := range tests {
		result := isJSON([]byte(test.content))
		if result != test.expected {
			t.Errorf("isJSON(%s) = %v; expected %v", test.content, result, test.expected)
		}
	}
}

// TestIsYAML tests the isYAML function
func TestIsYAML(t *testing.T) {
	tests := []struct {
		content  string
		expected bool
	}{
		{`name: test`, true},
		{`foo: bar
baz: qux`, true},
		{`{"name": "test"}`, false},
		{`<html></html>`, false},
	}

	for _, test := range tests {
		result := isYAML([]byte(test.content))
		if result != test.expected {
			t.Errorf("isYAML(%s) = %v; expected %v", test.content, result, test.expected)
		}
	}
}

// TestScrapeAPIDoc tests the ScrapeAPIDoc function with a mock server
func TestScrapeAPIDoc(t *testing.T) {
	// Create a test server that returns JSON
	jsonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"openapi": "3.0.0",
			"info": {
				"title": "Test API",
				"description": "API for testing",
				"version": "1.0.0"
			},
			"paths": {}
		}`))
	}))
	defer jsonServer.Close()

	// Create a test server that returns HTML
	htmlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
		<html>
		<head>
			<title>Test API Documentation</title>
			<meta name="description" content="API documentation for testing">
		</head>
		<body>
			<h1>Test API</h1>
			<p>This is a test API for demonstration purposes.</p>
			
			<h2>API Endpoints</h2>
			
			<h3>GET /users</h3>
			<p>Returns a list of users.</p>
		</body>
		</html>`))
	}))
	defer htmlServer.Close()

	// Test scraping JSON
	jsonDoc, err := ScrapeAPIDoc(jsonServer.URL + "/swagger")
	if err != nil {
		t.Fatalf("Failed to scrape JSON API doc: %v", err)
	}

	if jsonDoc.Title != "Test API" {
		t.Errorf("Expected title 'Test API', got '%s'", jsonDoc.Title)
	}

	// Test scraping HTML
	htmlDoc, err := ScrapeAPIDoc(htmlServer.URL + "/api/doc")
	if err != nil {
		t.Fatalf("Failed to scrape HTML API doc: %v", err)
	}

	if htmlDoc.Title != "Test API Documentation" {
		t.Errorf("Expected title 'Test API Documentation', got '%s'", htmlDoc.Title)
	}
}
