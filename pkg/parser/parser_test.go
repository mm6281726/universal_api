package parser

import (
	"testing"
	"universal_api/internal/models"
)

// Test data
const jsonTestData = `{
	"openapi": "3.0.0",
	"info": {
		"title": "Test API",
		"description": "API for testing",
		"version": "1.0.0"
	},
	"paths": {
		"/users": {
			"get": {
				"summary": "Get all users",
				"description": "Returns a list of users",
				"parameters": [
					{
						"name": "limit",
						"in": "query",
						"description": "Maximum number of users to return",
						"required": false,
						"schema": {
							"type": "integer"
						}
					}
				],
				"responses": {
					"200": {
						"description": "Successful operation"
					},
					"400": {
						"description": "Bad request"
					}
				}
			},
			"post": {
				"summary": "Create a user",
				"description": "Creates a new user",
				"parameters": [
					{
						"name": "name",
						"in": "body",
						"description": "User name",
						"required": true,
						"schema": {
							"type": "string"
						}
					}
				],
				"responses": {
					"201": {
						"description": "User created"
					},
					"400": {
						"description": "Bad request"
					}
				}
			}
		}
	}
}`

const yamlTestData = `openapi: 3.0.0
info:
  title: Test API
  description: API for testing
  version: 1.0.0
paths:
  /users:
    get:
      summary: Get all users
      description: Returns a list of users
      parameters:
        - name: limit
          in: query
          description: Maximum number of users to return
          required: false
          schema:
            type: integer
      responses:
        '200':
          description: Successful operation
        '400':
          description: Bad request
    post:
      summary: Create a user
      description: Creates a new user
      parameters:
        - name: name
          in: body
          description: User name
          required: true
          schema:
            type: string
      responses:
        '201':
          description: User created
        '400':
          description: Bad request
`

const htmlTestData = `<!DOCTYPE html>
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
    <h4>Parameters</h4>
    <table>
        <tr>
            <th>Name</th>
            <th>Description</th>
        </tr>
        <tr>
            <td>limit</td>
            <td>Maximum number of users to return</td>
        </tr>
    </table>
    <h4>Responses</h4>
    <ul>
        <li>200: Successful operation</li>
        <li>400: Bad request</li>
    </ul>
    
    <h3>POST /users</h3>
    <p>Creates a new user.</p>
    <h4>Parameters</h4>
    <table>
        <tr>
            <th>Name</th>
            <th>Description</th>
        </tr>
        <tr>
            <td>name</td>
            <td>User name</td>
        </tr>
    </table>
    <h4>Responses</h4>
    <ul>
        <li>201: User created</li>
        <li>400: Bad request</li>
    </ul>
</body>
</html>`

// TestJSONParser tests the JSON parser
func TestJSONParser(t *testing.T) {
	parser := &JSONParser{}
	
	apiDoc, err := parser.Parse([]byte(jsonTestData))
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	// Verify basic info
	if apiDoc.Title != "Test API" {
		t.Errorf("Expected title 'Test API', got '%s'", apiDoc.Title)
	}
	
	if apiDoc.Description != "API for testing" {
		t.Errorf("Expected description 'API for testing', got '%s'", apiDoc.Description)
	}
	
	if apiDoc.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", apiDoc.Version)
	}
	
	// Verify endpoints
	if len(apiDoc.Endpoints) != 2 {
		t.Fatalf("Expected 2 endpoints, got %d", len(apiDoc.Endpoints))
	}
	
	// Check first endpoint (GET /users)
	getEndpoint := findEndpoint(apiDoc.Endpoints, "GET", "/users")
	if getEndpoint == nil {
		t.Fatalf("GET /users endpoint not found")
	}
	
	if getEndpoint.Summary != "Get all users" {
		t.Errorf("Expected summary 'Get all users', got '%s'", getEndpoint.Summary)
	}
	
	if len(getEndpoint.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(getEndpoint.Parameters))
	}
	
	if len(getEndpoint.Responses) != 2 {
		t.Errorf("Expected 2 responses, got %d", len(getEndpoint.Responses))
	}
	
	// Check second endpoint (POST /users)
	postEndpoint := findEndpoint(apiDoc.Endpoints, "POST", "/users")
	if postEndpoint == nil {
		t.Fatalf("POST /users endpoint not found")
	}
	
	if postEndpoint.Summary != "Create a user" {
		t.Errorf("Expected summary 'Create a user', got '%s'", postEndpoint.Summary)
	}
}

// TestYAMLParser tests the YAML parser
func TestYAMLParser(t *testing.T) {
	parser := &YAMLParser{}
	
	apiDoc, err := parser.Parse([]byte(yamlTestData))
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}
	
	// Verify basic info
	if apiDoc.Title != "Test API" {
		t.Errorf("Expected title 'Test API', got '%s'", apiDoc.Title)
	}
	
	if apiDoc.Description != "API for testing" {
		t.Errorf("Expected description 'API for testing', got '%s'", apiDoc.Description)
	}
	
	if apiDoc.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", apiDoc.Version)
	}
	
	// Verify endpoints
	if len(apiDoc.Endpoints) != 2 {
		t.Fatalf("Expected 2 endpoints, got %d", len(apiDoc.Endpoints))
	}
	
	// Check first endpoint (GET /users)
	getEndpoint := findEndpoint(apiDoc.Endpoints, "GET", "/users")
	if getEndpoint == nil {
		t.Fatalf("GET /users endpoint not found")
	}
	
	if getEndpoint.Summary != "Get all users" {
		t.Errorf("Expected summary 'Get all users', got '%s'", getEndpoint.Summary)
	}
}

// TestHTMLParser tests the HTML parser
func TestHTMLParser(t *testing.T) {
	parser := &HTMLParser{}
	
	apiDoc, err := parser.Parse([]byte(htmlTestData))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}
	
	// Verify basic info
	if apiDoc.Title != "Test API Documentation" {
		t.Errorf("Expected title 'Test API Documentation', got '%s'", apiDoc.Title)
	}
	
	if apiDoc.Description != "API documentation for testing" {
		t.Errorf("Expected description 'API documentation for testing', got '%s'", apiDoc.Description)
	}
	
	// Verify endpoints
	if len(apiDoc.Endpoints) < 2 {
		t.Fatalf("Expected at least 2 endpoints, got %d", len(apiDoc.Endpoints))
	}
	
	// Check if GET /users endpoint exists
	getEndpoint := findEndpoint(apiDoc.Endpoints, "GET", "/users")
	if getEndpoint == nil {
		t.Fatalf("GET /users endpoint not found")
	}
	
	// Check if POST /users endpoint exists
	postEndpoint := findEndpoint(apiDoc.Endpoints, "POST", "/users")
	if postEndpoint == nil {
		t.Fatalf("POST /users endpoint not found")
	}
}

// Helper function to find an endpoint by method and path
func findEndpoint(endpoints []models.Endpoint, method, path string) *models.Endpoint {
	for i, endpoint := range endpoints {
		if endpoint.Method == method && endpoint.Path == path {
			return &endpoints[i]
		}
	}
	return nil
}
