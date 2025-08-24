package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// WebHandlers handles web page endpoints
type WebHandlers struct {
	*BaseHandler
}

// NewWebHandlers creates a new web handlers instance
func NewWebHandlers(base *BaseHandler) *WebHandlers {
	return &WebHandlers{BaseHandler: base}
}

// ServeHomePage serves the main blog page
func (h *WebHandlers) ServeHomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get posts from Typesense instead of MySQL for public pages
	posts, err := h.SearchService.GetAllPostsFromIndex(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := generateHomePageHTML(posts)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// ServePostDetail serves the post detail page
func (h *WebHandlers) ServePostDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path using gorilla/mux
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.GetPost(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	html := generatePostDetailHTML(post)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
