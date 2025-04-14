package parser

import (
	"errors"
	"universal_api/internal/models"
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

// Parse implements the Parser interface for JSON
func (p *JSONParser) Parse(content []byte) (*models.APIDoc, error) {
	// This would be implemented to parse JSON API documentation
	// For now, return an error
	return nil, errors.New("JSON parser not implemented yet")
}

// YAMLParser parses YAML API documentation
type YAMLParser struct{}

// Parse implements the Parser interface for YAML
func (p *YAMLParser) Parse(content []byte) (*models.APIDoc, error) {
	// This would be implemented to parse YAML API documentation
	// For now, return an error
	return nil, errors.New("YAML parser not implemented yet")
}

// HTMLParser parses HTML API documentation
type HTMLParser struct{}

// Parse implements the Parser interface for HTML
func (p *HTMLParser) Parse(content []byte) (*models.APIDoc, error) {
	// This would be implemented to parse HTML API documentation
	// For now, return an error
	return nil, errors.New("HTML parser not implemented yet")
}
