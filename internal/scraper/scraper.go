package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"universal_api/internal/models"
	"universal_api/pkg/parser"
)

// ScrapeAPIDoc scrapes API documentation from the given URL
func ScrapeAPIDoc(url string) (*models.APIDoc, error) {
	// Check if the URL is for a known API documentation format
	if isSwaggerURL(url) {
		return scrapeSwaggerDoc(url)
	} else if isRESTDocURL(url) {
		return scrapeGenericRESTDoc(url)
	}

	// Default to generic scraping
	return scrapeGenericDoc(url)
}

// isSwaggerURL checks if the URL is for Swagger/OpenAPI documentation
func isSwaggerURL(url string) bool {
	return strings.Contains(url, "swagger") ||
		strings.Contains(url, "openapi") ||
		strings.Contains(url, "api-docs")
}

// isRESTDocURL checks if the URL is for REST API documentation
func isRESTDocURL(url string) bool {
	return strings.Contains(url, "api") &&
		(strings.Contains(url, "doc") || strings.Contains(url, "reference"))
}

// scrapeSwaggerDoc scrapes Swagger/OpenAPI documentation
func scrapeSwaggerDoc(url string) (*models.APIDoc, error) {
	// Make HTTP request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Determine content type
	contentType := resp.Header.Get("Content-Type")

	// Create parser based on content type
	var p parser.Parser
	if strings.Contains(contentType, "json") {
		p = &parser.JSONParser{}
	} else if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "yml") {
		p = &parser.YAMLParser{}
	} else {
		// Try to detect format from content
		if isJSON(content) {
			p = &parser.JSONParser{}
		} else if isYAML(content) {
			p = &parser.YAMLParser{}
		} else {
			// Default to JSON parser
			p = &parser.JSONParser{}
		}
	}

	// Parse the content
	apiDoc, err := p.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Swagger/OpenAPI documentation: %w", err)
	}

	// Set URL and timestamps
	apiDoc.URL = url
	apiDoc.CreatedAt = time.Now()
	apiDoc.UpdatedAt = time.Now()

	return apiDoc, nil
}

// scrapeGenericRESTDoc scrapes generic REST API documentation
func scrapeGenericRESTDoc(url string) (*models.APIDoc, error) {
	// Make HTTP request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Determine content type
	contentType := resp.Header.Get("Content-Type")

	// Create parser based on content type
	var p parser.Parser
	if strings.Contains(contentType, "html") {
		p = &parser.HTMLParser{}
	} else if strings.Contains(contentType, "json") {
		p = &parser.JSONParser{}
	} else if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "yml") {
		p = &parser.YAMLParser{}
	} else {
		// Default to HTML parser for REST docs
		p = &parser.HTMLParser{}
	}

	// Parse the content
	apiDoc, err := p.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REST API documentation: %w", err)
	}

	// Set URL and timestamps
	apiDoc.URL = url
	apiDoc.CreatedAt = time.Now()
	apiDoc.UpdatedAt = time.Now()

	return apiDoc, nil
}

// scrapeGenericDoc scrapes generic API documentation
func scrapeGenericDoc(url string) (*models.APIDoc, error) {
	// Make HTTP request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Determine content type
	contentType := resp.Header.Get("Content-Type")

	// Create parser based on content type
	var p parser.Parser
	if strings.Contains(contentType, "html") {
		p = &parser.HTMLParser{}
	} else if strings.Contains(contentType, "json") {
		p = &parser.JSONParser{}
	} else if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "yml") {
		p = &parser.YAMLParser{}
	} else {
		// Try to detect format from content
		if isJSON(content) {
			p = &parser.JSONParser{}
		} else if isYAML(content) {
			p = &parser.YAMLParser{}
		} else {
			// Default to HTML parser
			p = &parser.HTMLParser{}
		}
	}

	// Parse the content
	apiDoc, err := p.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API documentation: %w", err)
	}

	// Set URL and timestamps
	apiDoc.URL = url
	apiDoc.CreatedAt = time.Now()
	apiDoc.UpdatedAt = time.Now()

	return apiDoc, nil
}

// Helper functions

// isJSON checks if content is likely JSON
func isJSON(content []byte) bool {
	trimmed := strings.TrimSpace(string(content))
	return (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]"))
}

// isYAML checks if content is likely YAML
func isYAML(content []byte) bool {
	// Simple heuristic: YAML often has key: value pairs
	trimmed := strings.TrimSpace(string(content))

	// If it looks like JSON, it's not YAML
	if isJSON(content) {
		return false
	}

	lines := strings.Split(trimmed, "\n")
	for _, line := range lines {
		// Look for key: value pattern but exclude JSON-like patterns
		if strings.Contains(line, ":") &&
			!strings.Contains(line, "{:") &&
			!strings.Contains(line, "\":") {
			return true
		}
	}
	return false
}
