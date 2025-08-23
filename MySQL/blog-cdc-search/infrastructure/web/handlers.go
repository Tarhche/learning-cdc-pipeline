package web

import (
	"net/http"

	"blog-cdc-search/application/service"
	"blog-cdc-search/infrastructure/web/handlers"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	API       *handlers.APIHandlers
	Web       *handlers.WebHandlers
	Dashboard *handlers.DashboardHandlers
	Search    *handlers.SearchHandlers
}

// NewHandlers creates a new Handlers instance
func NewHandlers(postService *service.PostService, searchService handlers.SearchService) *Handlers {
	base := handlers.NewBaseHandler(postService, searchService)
	return &Handlers{
		API:       handlers.NewAPIHandlers(base),
		Web:       handlers.NewWebHandlers(base),
		Dashboard: handlers.NewDashboardHandlers(base),
		Search:    handlers.NewSearchHandlers(base),
	}
}

// Legacy methods for backward compatibility - delegate to appropriate handlers

func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	h.API.CreatePost(w, r)
}

func (h *Handlers) GetPost(w http.ResponseWriter, r *http.Request) {
	h.API.GetPost(w, r)
}

func (h *Handlers) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	h.API.GetAllPosts(w, r)
}

func (h *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	h.API.UpdatePost(w, r)
}

func (h *Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	h.API.DeletePost(w, r)
}

func (h *Handlers) ServeHomePage(w http.ResponseWriter, r *http.Request) {
	h.Web.ServeHomePage(w, r)
}

func (h *Handlers) ServePostDetail(w http.ResponseWriter, r *http.Request) {
	h.Web.ServePostDetail(w, r)
}

func (h *Handlers) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	h.Dashboard.ServeDashboard(w, r)
}

func (h *Handlers) ServeCreateForm(w http.ResponseWriter, r *http.Request) {
	h.Dashboard.ServeCreateForm(w, r)
}

func (h *Handlers) ServeEditForm(w http.ResponseWriter, r *http.Request) {
	h.Dashboard.ServeEditForm(w, r)
}

// Search methods
func (h *Handlers) SearchPosts(w http.ResponseWriter, r *http.Request) {
	h.Search.SearchPosts(w, r)
}

func (h *Handlers) SearchPostsGet(w http.ResponseWriter, r *http.Request) {
	h.Search.SearchPostsGet(w, r)
}
