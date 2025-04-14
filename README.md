# Universal API

A RESTful API service written in Go that scrapes and unifies API documentation from different providers.

## Overview

Universal API allows you to submit URLs to API documentation, which will then be scraped and saved in a structured format. This makes it easier to work with multiple APIs by providing a consistent way to access their documentation.

## Features

- Submit URLs to API documentation for scraping
- Automatic detection of API documentation format (Swagger/OpenAPI, REST, etc.)
- Unified storage of API documentation
- RESTful API for accessing the scraped API data

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Running the Application

1. Clone the repository
2. Run the application using Docker Compose:

```bash
docker-compose up
```

3. The API will be available at http://localhost:8081

## API Endpoints

### Submit API Documentation

```
POST /api/v1/docs
```

Request body:
```json
{
  "url": "https://example.com/api-docs",
  "description": "Example API Documentation"
}
```

### Get All API Docs

```
GET /api/v1/docs
```

### Get API Doc by ID

```
GET /api/v1/docs/:id
```

## Project Structure

- `cmd/api`: Main application entry point
- `internal/models`: Data models
- `internal/scraper`: API documentation scraper
- `internal/storage`: Storage layer
- `pkg/parser`: Parsers for different API documentation formats

## License

This project is licensed under the MIT License - see the LICENSE file for details.
