package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"universal_api/internal/models"
	"universal_api/internal/scraper"
	"universal_api/internal/storage"
)

// Global storage instance
var store storage.Storage

func main() {
	// Initialize storage
	store = storage.NewMemoryStorage()

	r := gin.Default()

	// Setup routes
	setupRoutes(r)

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Submit a new API documentation URL for scraping
		api.POST("/docs", submitAPIDoc)

		// Get all API docs
		api.GET("/docs", getAllAPIDocs)

		// Get a specific API doc by ID
		api.GET("/docs/:id", getAPIDocByID)
	}
}

// Handler to submit a new API documentation URL
func submitAPIDoc(c *gin.Context) {
	var request models.APIDocRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate URL
	if request.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// Scrape the API documentation
	apiDoc, err := scraper.ScrapeAPIDoc(request.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scrape API documentation: " + err.Error()})
		return
	}

	// Set description from request if provided
	if request.Description != "" {
		apiDoc.Description = request.Description
	}

	// Save the API doc
	if err := store.SaveAPIDoc(apiDoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save API documentation: " + err.Error()})
		return
	}

	// Return the API doc
	c.JSON(http.StatusOK, apiDoc)
}

// Handler to get all API docs
func getAllAPIDocs(c *gin.Context) {
	docs, err := store.GetAllAPIDocs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API docs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, docs)
}

// Handler to get a specific API doc by ID
func getAPIDocByID(c *gin.Context) {
	id := c.Param("id")

	doc, err := store.GetAPIDoc(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API doc not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}
