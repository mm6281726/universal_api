package ui

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"universal_api/internal/scraper"
	"universal_api/internal/storage"
)

// Handler handles UI requests
type Handler struct {
	templates *template.Template
	store     storage.Storage
	limiter   *RateLimiter
}

// NewHandler creates a new UI handler
func NewHandler(store storage.Storage) *Handler {
	// Parse templates
	templates := template.New("").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	})

	templatePath := filepath.Join("internal", "ui", "templates")
	templates, err := templates.ParseGlob(filepath.Join(templatePath, "*.html"))
	if err != nil {
		panic(err)
	}

	return &Handler{
		templates: templates,
		store:     store,
		limiter:   NewRateLimiter(1, 5), // 1 request per domain every 5 seconds
	}
}

// RegisterRoutes registers UI routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/docs", h.handleDocsList)
	mux.HandleFunc("/docs/", h.handleDocDetail)
	mux.HandleFunc("/scrape", h.handleScrape)
}

// handleIndex handles the index page
func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get the most recent API docs (up to 5)
	docs, err := h.store.GetAllAPIDocs()
	if err != nil {
		h.renderError(w, "Failed to get API docs: "+err.Error())
		return
	}

	// Limit to 5 most recent docs
	recentDocs := docs
	if len(docs) > 5 {
		recentDocs = docs[len(docs)-5:]
	}

	data := map[string]interface{}{
		"Title":   "Home",
		"APIDocs": recentDocs,
	}

	h.renderTemplate(w, "index", data)
}

// handleDocsList handles the docs list page
func (h *Handler) handleDocsList(w http.ResponseWriter, r *http.Request) {
	docs, err := h.store.GetAllAPIDocs()
	if err != nil {
		h.renderError(w, "Failed to get API docs: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"Title":   "API Documentation",
		"APIDocs": docs,
	}

	h.renderTemplate(w, "docs_list", data)
}

// handleDocDetail handles the doc detail page
func (h *Handler) handleDocDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/docs/")
	if id == "" {
		http.Redirect(w, r, "/docs", http.StatusSeeOther)
		return
	}

	doc, err := h.store.GetAPIDoc(id)
	if err != nil {
		h.renderError(w, "API doc not found: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"Title":  doc.Title,
		"APIDoc": doc,
	}

	h.renderTemplate(w, "doc_detail", data)
}

// handleScrape handles the scrape action
func (h *Handler) handleScrape(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		h.renderError(w, "URL is required")
		return
	}

	// Check rate limit
	if !h.limiter.Allow(url) {
		h.renderError(w, "Rate limit exceeded for this domain. Please try again later.")
		return
	}

	// Scrape the API documentation
	apiDoc, err := scraper.ScrapeAPIDoc(url)
	if err != nil {
		h.renderError(w, "Failed to scrape API documentation: "+err.Error())
		return
	}

	// Save the API doc
	if err := h.store.SaveAPIDoc(apiDoc); err != nil {
		h.renderError(w, "Failed to save API documentation: "+err.Error())
		return
	}

	// Redirect to the doc detail page
	http.Redirect(w, r, "/docs/"+apiDoc.ID, http.StatusSeeOther)
}

// renderTemplate renders a template with the given data
func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.templates.ExecuteTemplate(w, name+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderError renders an error page
func (h *Handler) renderError(w http.ResponseWriter, message string) {
	data := map[string]interface{}{
		"Title": "Error",
		"Error": message,
	}

	h.renderTemplate(w, "error", data)
}
