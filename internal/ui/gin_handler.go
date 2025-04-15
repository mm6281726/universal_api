package ui

import (
	"html/template"
	"net/http"
	"strings"

	"universal_api/internal/scraper"
	"universal_api/internal/storage"

	"github.com/gin-gonic/gin"
)

// GinHandler handles UI requests for Gin
type GinHandler struct {
	store   storage.Storage
	limiter *RateLimiter
}

// NewGinHandler creates a new Gin UI handler
func NewGinHandler(store storage.Storage) *GinHandler {
	return &GinHandler{
		store:   store,
		limiter: NewRateLimiter(1, 5), // 1 request per domain every 5 seconds
	}
}

// RegisterRoutes registers UI routes with Gin
func (h *GinHandler) RegisterRoutes(r *gin.Engine) {
	// Serve static files
	r.Static("/static", "./internal/ui/static")

	// Add template functions
	r.SetFuncMap(template.FuncMap{
		"lower": strings.ToLower,
	})

	// Load HTML templates
	r.LoadHTMLGlob("internal/ui/templates/*")

	// UI routes
	r.GET("/", h.handleIndex)
	r.GET("/docs", h.handleDocsList)
	r.GET("/docs/:id", h.handleDocDetail)
	r.POST("/scrape", h.handleScrape)
}

// handleIndex handles the index page
func (h *GinHandler) handleIndex(c *gin.Context) {
	// Get the most recent API docs (up to 5)
	docs, err := h.store.GetAllAPIDocs()
	if err != nil {
		h.renderError(c, "Failed to get API docs: "+err.Error())
		return
	}

	// Limit to 5 most recent docs
	recentDocs := docs
	if len(docs) > 5 {
		recentDocs = docs[len(docs)-5:]
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":   "Home",
		"APIDocs": recentDocs,
	})
}

// handleDocsList handles the docs list page
func (h *GinHandler) handleDocsList(c *gin.Context) {
	docs, err := h.store.GetAllAPIDocs()
	if err != nil {
		h.renderError(c, "Failed to get API docs: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "docs_list.tmpl", gin.H{
		"Title":   "API Documentation",
		"APIDocs": docs,
	})
}

// handleDocDetail handles the doc detail page
func (h *GinHandler) handleDocDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Redirect(http.StatusSeeOther, "/docs")
		return
	}

	doc, err := h.store.GetAPIDoc(id)
	if err != nil {
		h.renderError(c, "API doc not found: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "doc_detail.tmpl", gin.H{
		"Title":  doc.Title,
		"APIDoc": doc,
	})
}

// handleScrape handles the scrape action
func (h *GinHandler) handleScrape(c *gin.Context) {
	url := c.PostForm("url")

	if url == "" {
		h.renderError(c, "URL is required")
		return
	}

	// Check rate limit
	if !h.limiter.Allow(url) {
		h.renderError(c, "Rate limit exceeded for this domain. Please try again later.")
		return
	}

	// Scrape the API documentation
	apiDoc, err := scraper.ScrapeAPIDoc(url)
	if err != nil {
		h.renderError(c, "Failed to scrape API documentation: "+err.Error())
		return
	}

	// Save the API doc
	if err := h.store.SaveAPIDoc(apiDoc); err != nil {
		h.renderError(c, "Failed to save API documentation: "+err.Error())
		return
	}

	// Redirect to the doc detail page
	c.Redirect(http.StatusSeeOther, "/docs/"+apiDoc.ID)
}

// renderError renders an error page
func (h *GinHandler) renderError(c *gin.Context, message string) {
	c.HTML(http.StatusOK, "error.tmpl", gin.H{
		"Title": "Error",
		"Error": message,
	})
}
