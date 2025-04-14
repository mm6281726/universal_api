package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"universal_api/internal/models"
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
	// This would be implemented to parse Swagger/OpenAPI JSON or YAML
	// For now, return a placeholder
	return &models.APIDoc{
		ID:          fmt.Sprintf("swagger-%d", time.Now().Unix()),
		URL:         url,
		Title:       "Swagger API Documentation",
		Description: "Scraped from Swagger/OpenAPI documentation",
		Version:     "1.0.0",
		Endpoints:   []models.Endpoint{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// scrapeGenericRESTDoc scrapes generic REST API documentation
func scrapeGenericRESTDoc(url string) (*models.APIDoc, error) {
	// This would be implemented to parse generic REST API documentation
	// For now, return a placeholder
	return &models.APIDoc{
		ID:          fmt.Sprintf("rest-%d", time.Now().Unix()),
		URL:         url,
		Title:       "REST API Documentation",
		Description: "Scraped from generic REST API documentation",
		Version:     "1.0.0",
		Endpoints:   []models.Endpoint{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
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
		return nil, errors.New(fmt.Sprintf("HTTP request failed with status code: %d", resp.StatusCode))
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract title
	title := doc.Find("title").Text()
	if title == "" {
		title = "Unknown API"
	}

	// Extract description
	description := ""
	doc.Find("meta[name=description]").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists {
			description = content
		}
	})

	// Create API doc
	apiDoc := &models.APIDoc{
		ID:          fmt.Sprintf("generic-%d", time.Now().Unix()),
		URL:         url,
		Title:       title,
		Description: description,
		Version:     "Unknown",
		Endpoints:   []models.Endpoint{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Try to extract endpoints
	// This is a very basic implementation and would need to be enhanced
	doc.Find("h2, h3").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(strings.ToLower(text), "api") || 
		   strings.Contains(strings.ToLower(text), "endpoint") {
			// Found a potential endpoint section
			endpoint := models.Endpoint{
				Path:        "Unknown",
				Method:      "Unknown",
				Summary:     text,
				Description: s.Next().Text(),
				Parameters:  []models.Parameter{},
				Responses:   []models.Response{},
			}
			apiDoc.Endpoints = append(apiDoc.Endpoints, endpoint)
		}
	})

	return apiDoc, nil
}
