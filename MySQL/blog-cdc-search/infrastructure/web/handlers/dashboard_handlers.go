package handlers

import (
	"net/http"
	"strconv"
)

// DashboardHandlers handles admin dashboard endpoints
type DashboardHandlers struct {
	*BaseHandler
}

// NewDashboardHandlers creates a new dashboard handlers instance
func NewDashboardHandlers(base *BaseHandler) *DashboardHandlers {
	return &DashboardHandlers{BaseHandler: base}
}

// ServeDashboard serves the admin dashboard page
func (h *DashboardHandlers) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.PostService.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := generateDashboardHTML(posts)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// ServeCreateForm serves the create post form
func (h *DashboardHandlers) ServeCreateForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	html := generateCreateFormHTML()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// ServeEditForm serves the edit post form
func (h *DashboardHandlers) ServeEditForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
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

	html := generateEditFormHTML(post)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
