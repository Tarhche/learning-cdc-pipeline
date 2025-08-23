package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"blog-cdc-search/application/service"
	"blog-cdc-search/domain"
)

// SearchServiceAdapter adapts the actual SearchService to the interface
type SearchServiceAdapter struct {
	service *service.SearchService
}

func (a *SearchServiceAdapter) SearchPosts(ctx context.Context, params interface{}) (interface{}, error) {
	// Convert interface{} to SearchParams
	searchParams, ok := params.(service.SearchParams)
	if !ok {
		// Try to convert from map[string]interface{} if that's what was passed
		if paramMap, ok := params.(map[string]interface{}); ok {
			searchParams = service.SearchParams{
				Query:    paramMap["query"].(string),
				Page:     paramMap["page"].(int),
				PerPage:  paramMap["per_page"].(int),
				SortBy:   paramMap["sort_by"].(string),
				FilterBy: paramMap["filter_by"].(string),
			}
		} else {
			// Return error for unsupported type
			return nil, fmt.Errorf("unsupported params type: %T", params)
		}
	}

	return a.service.SearchPosts(ctx, searchParams)
}

func (a *SearchServiceAdapter) GetAllPostsFromIndex(ctx context.Context) ([]*domain.Post, error) {
	return a.service.GetAllPostsFromIndex(ctx)
}

// SearchHandlers handles search-related endpoints
type SearchHandlers struct {
	*BaseHandler
}

// NewSearchHandlers creates a new search handlers instance
func NewSearchHandlers(base *BaseHandler) *SearchHandlers {
	return &SearchHandlers{BaseHandler: base}
}

// SearchPosts handles POST /api/search
func (h *SearchHandlers) SearchPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query    string `json:"query"`
		Page     int    `json:"page"`
		PerPage  int    `json:"per_page"`
		SortBy   string `json:"sort_by"`
		FilterBy string `json:"filter_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create search parameters
	searchParams := service.SearchParams{
		Query:    req.Query,
		Page:     req.Page,
		PerPage:  req.PerPage,
		SortBy:   req.SortBy,
		FilterBy: req.FilterBy,
	}

	// Perform search
	results, err := h.SearchService.SearchPosts(r.Context(), searchParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// SearchPostsGet handles GET /api/search for simple query-based searches
func (h *SearchHandlers) SearchPostsGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPage := 10
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	// Create search parameters
	searchParams := service.SearchParams{
		Query:   query,
		Page:    page,
		PerPage: perPage,
	}

	// Perform search
	results, err := h.SearchService.SearchPosts(r.Context(), searchParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
