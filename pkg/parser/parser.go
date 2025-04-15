package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"universal_api/internal/models"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v3"
)

// Parser interface for different API documentation formats
type Parser interface {
	Parse(content []byte) (*models.APIDoc, error)
}

// ParserFactory creates a parser based on the content type
func ParserFactory(contentType string) (Parser, error) {
	switch contentType {
	case "application/json", "text/json":
		return &JSONParser{}, nil
	case "application/yaml", "text/yaml", "application/x-yaml":
		return &YAMLParser{}, nil
	case "text/html":
		return &HTMLParser{}, nil
	default:
		return nil, errors.New("unsupported content type")
	}
}

// JSONParser parses JSON API documentation
type JSONParser struct{}

// OpenAPIDoc represents a simplified OpenAPI/Swagger document structure
type OpenAPIDoc struct {
	Openapi     string                 `json:"openapi,omitempty"`
	Swagger     string                 `json:"swagger,omitempty"`
	Info        OpenAPIInfo            `json:"info"`
	Paths       map[string]PathItem    `json:"paths"`
	Components  *OpenAPIComponents     `json:"components,omitempty"`
	Definitions map[string]interface{} `json:"definitions,omitempty"` // For Swagger 2.0
}

// OpenAPIInfo contains metadata about the API
type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

// OpenAPIComponents contains reusable objects for different aspects of the OAS
type OpenAPIComponents struct {
	Schemas map[string]interface{} `json:"schemas,omitempty"`
}

// PathItem describes the operations available on a single path
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
}

// Operation describes a single API operation on a path
type Operation struct {
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	OperationID string                 `json:"operationId,omitempty"`
	Parameters  []Parameter            `json:"parameters,omitempty"`
	Responses   map[string]interface{} `json:"responses,omitempty"`
}

// Parameter describes a single operation parameter
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // query, path, header, cookie, body
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
	Type        string  `json:"type,omitempty"` // For Swagger 2.0
}

// Schema represents a Schema Object in OpenAPI
type Schema struct {
	Type string `json:"type,omitempty"`
}

// Parse implements the Parser interface for JSON
func (p *JSONParser) Parse(content []byte) (*models.APIDoc, error) {
	// Try to parse as OpenAPI/Swagger
	var openAPIDoc OpenAPIDoc
	if err := json.Unmarshal(content, &openAPIDoc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON as OpenAPI: %w", err)
	}

	// Validate that it's actually an OpenAPI document
	if openAPIDoc.OpenAPI() == "" {
		return nil, errors.New("JSON does not appear to be an OpenAPI/Swagger document")
	}

	// Create API doc
	apiDoc := &models.APIDoc{
		ID:          fmt.Sprintf("openapi-%d", time.Now().Unix()),
		Title:       openAPIDoc.Info.Title,
		Description: openAPIDoc.Info.Description,
		Version:     openAPIDoc.Info.Version,
		Endpoints:   []models.Endpoint{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Extract endpoints
	for path, pathItem := range openAPIDoc.Paths {
		// Process each HTTP method
		operations := pathItem.Operations()
		for method, operation := range operations {
			// Create endpoint
			endpoint := models.Endpoint{
				Path:        path,
				Method:      method,
				Summary:     operation.Summary,
				Description: operation.Description,
				Parameters:  []models.Parameter{},
				Responses:   []models.Response{},
			}

			// Add parameters
			for _, param := range operation.Parameters {
				paramType := param.Type
				if param.Schema != nil && param.Schema.Type != "" {
					paramType = param.Schema.Type
				}

				endpoint.Parameters = append(endpoint.Parameters, models.Parameter{
					Name:        param.Name,
					In:          param.In,
					Required:    param.Required,
					Type:        paramType,
					Description: param.Description,
				})
			}

			// Add responses
			for statusCode, responseObj := range operation.Responses {
				// Try to extract description from response object
				description := ""
				if respMap, ok := responseObj.(map[string]interface{}); ok {
					if desc, ok := respMap["description"].(string); ok {
						description = desc
					}
				}

				// Convert status code to int
				code := 0
				if statusCode == "default" {
					code = 0
				} else {
					fmt.Sscanf(statusCode, "%d", &code)
				}

				endpoint.Responses = append(endpoint.Responses, models.Response{
					StatusCode:  code,
					Description: description,
					// Schema is omitted for simplicity
				})
			}

			apiDoc.Endpoints = append(apiDoc.Endpoints, endpoint)
		}
	}

	return apiDoc, nil
}

// OpenAPI returns the OpenAPI version (either from openapi or swagger field)
func (doc *OpenAPIDoc) OpenAPI() string {
	if doc.Openapi != "" {
		return doc.Openapi
	}
	return doc.Swagger
}

// Operations returns a map of HTTP method to Operation for a PathItem
func (item *PathItem) Operations() map[string]Operation {
	result := make(map[string]Operation)

	if item.Get != nil {
		result["GET"] = *item.Get
	}
	if item.Post != nil {
		result["POST"] = *item.Post
	}
	if item.Put != nil {
		result["PUT"] = *item.Put
	}
	if item.Delete != nil {
		result["DELETE"] = *item.Delete
	}
	if item.Options != nil {
		result["OPTIONS"] = *item.Options
	}
	if item.Head != nil {
		result["HEAD"] = *item.Head
	}
	if item.Patch != nil {
		result["PATCH"] = *item.Patch
	}

	return result
}

// YAMLParser parses YAML API documentation
type YAMLParser struct{}

// Parse implements the Parser interface for YAML
func (p *YAMLParser) Parse(content []byte) (*models.APIDoc, error) {
	// Convert YAML to JSON
	var yamlObj interface{}
	if err := yaml.Unmarshal(content, &yamlObj); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Convert YAML object to JSON
	jsonData, err := json.Marshal(yamlObj)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}

	// Use the JSON parser to parse the converted data
	jsonParser := &JSONParser{}
	return jsonParser.Parse(jsonData)
}

// HTMLParser parses HTML API documentation
type HTMLParser struct{}

// Parse implements the Parser interface for HTML
func (p *HTMLParser) Parse(content []byte) (*models.APIDoc, error) {
	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(content)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
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

	// If no meta description, try to find a description in the content
	if description == "" {
		// Look for a description in the first paragraph or div
		description = doc.Find("p").First().Text()
		if description == "" {
			description = doc.Find("div").First().Text()
		}
		// Truncate if too long
		if len(description) > 200 {
			description = description[:197] + "..."
		}
	}

	// Create API doc
	apiDoc := &models.APIDoc{
		ID:          fmt.Sprintf("html-%d", time.Now().Unix()),
		Title:       title,
		Description: description,
		Version:     "Unknown",
		Endpoints:   []models.Endpoint{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Extract endpoints
	// Look for common patterns in API documentation

	// Method 1: Look for headings that might indicate endpoints
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		text = strings.TrimSpace(text)

		// Skip if empty
		if text == "" {
			return
		}

		// Check if this heading looks like an API endpoint
		if containsEndpointIndicators(text) {
			// Extract method and path
			method, path := extractMethodAndPath(text)

			// Get description from the next element
			description := ""
			next := s.Next()
			if next.Length() > 0 {
				description = next.Text()
			}

			// Create endpoint
			endpoint := models.Endpoint{
				Path:        path,
				Method:      method,
				Summary:     text,
				Description: description,
				Parameters:  []models.Parameter{},
				Responses:   []models.Response{},
			}

			// Look for parameters in tables or lists
			paramSection := s.NextUntil("h1, h2, h3, h4, h5, h6")
			paramSection.Find("table tr").Each(func(i int, row *goquery.Selection) {
				// Skip header row
				if i == 0 {
					return
				}

				// Extract parameter info from table row
				cells := row.Find("td")
				if cells.Length() >= 2 {
					name := cells.Eq(0).Text()
					desc := cells.Eq(1).Text()

					// Try to determine parameter type and location
					paramType := "string"
					paramIn := "query"
					required := false

					// Check if path parameter
					if strings.Contains(path, "{"+name+"}") || strings.Contains(path, ":"+name) {
						paramIn = "path"
						required = true
					}

					// Add parameter
					endpoint.Parameters = append(endpoint.Parameters, models.Parameter{
						Name:        name,
						In:          paramIn,
						Required:    required,
						Type:        paramType,
						Description: desc,
					})
				}
			})

			// Look for response codes
			paramSection.Find("code, .code, pre").Each(func(i int, code *goquery.Selection) {
				text := code.Text()

				// Look for HTTP status codes
				statusMatches := []string{"200", "201", "400", "401", "403", "404", "500"}
				for _, status := range statusMatches {
					if strings.Contains(text, status) {
						statusCode := 0
						fmt.Sscanf(status, "%d", &statusCode)

						// Check if this status code is already added
						alreadyAdded := false
						for _, resp := range endpoint.Responses {
							if resp.StatusCode == statusCode {
								alreadyAdded = true
								break
							}
						}

						if !alreadyAdded {
							endpoint.Responses = append(endpoint.Responses, models.Response{
								StatusCode:  statusCode,
								Description: getStatusCodeDescription(statusCode),
							})
						}
					}
				}
			})

			apiDoc.Endpoints = append(apiDoc.Endpoints, endpoint)
		}
	})

	// Method 2: Look for code blocks that might contain API endpoints
	doc.Find("pre, code, .code").Each(func(i int, s *goquery.Selection) {
		text := s.Text()

		// Check for common API request patterns
		for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH"} {
			if strings.Contains(text, method+" ") || strings.Contains(text, method+"\t") {
				lines := strings.Split(text, "\n")
				for _, line := range lines {
					if strings.Contains(line, method) {
						// Extract path
						parts := strings.Fields(line)
						if len(parts) >= 2 {
							path := parts[1]

							// Clean up path
							path = strings.TrimPrefix(path, "http://")
							path = strings.TrimPrefix(path, "https://")
							if idx := strings.Index(path, "/"); idx > 0 {
								path = path[idx:]
							}

							// Create endpoint
							endpoint := models.Endpoint{
								Path:        path,
								Method:      method,
								Summary:     line,
								Description: "",
								Parameters:  []models.Parameter{},
								Responses:   []models.Response{},
							}

							// Add a default response
							endpoint.Responses = append(endpoint.Responses, models.Response{
								StatusCode:  200,
								Description: "OK",
							})

							// Check if this endpoint is already added
							alreadyAdded := false
							for _, ep := range apiDoc.Endpoints {
								if ep.Path == path && ep.Method == method {
									alreadyAdded = true
									break
								}
							}

							if !alreadyAdded {
								apiDoc.Endpoints = append(apiDoc.Endpoints, endpoint)
							}
						}
					}
				}
			}
		}
	})

	return apiDoc, nil
}

// Helper functions for HTML parser

// containsEndpointIndicators checks if text contains indicators of an API endpoint
func containsEndpointIndicators(text string) bool {
	lowerText := strings.ToLower(text)

	// Check for common endpoint indicators
	if strings.Contains(lowerText, "api") ||
		strings.Contains(lowerText, "endpoint") ||
		strings.Contains(lowerText, "route") ||
		strings.Contains(lowerText, "request") {
		return true
	}

	// Check for HTTP methods
	for _, method := range []string{"get", "post", "put", "delete", "patch"} {
		if strings.Contains(lowerText, method) {
			return true
		}
	}

	// Check for URL path patterns
	if strings.Contains(text, "/") &&
		(strings.Contains(text, "{") || strings.Contains(text, ":")) {
		return true
	}

	return false
}

// extractMethodAndPath extracts HTTP method and path from text
func extractMethodAndPath(text string) (string, string) {
	// Default values
	method := "GET"
	path := "Unknown"

	// Check for HTTP methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	for _, m := range methods {
		if strings.Contains(strings.ToUpper(text), m) {
			method = m
			break
		}
	}

	// Extract path - look for patterns like /path, /path/{param}, etc.
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "/") {
			path = word
			break
		}
	}

	return method, path
}

// getStatusCodeDescription returns a description for common HTTP status codes
func getStatusCodeDescription(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown Status Code"
	}
}
