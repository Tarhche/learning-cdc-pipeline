package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-cdc-search/application/service"
	"blog-cdc-search/domain"
)

// MockSearchService is a mock implementation for testing
type MockSearchService struct {
	searchResults *service.SearchResponse
	searchError   error
}

func (m *MockSearchService) SearchPosts(ctx context.Context, params interface{}) (interface{}, error) {
	if m.searchError != nil {
		return nil, m.searchError
	}
	return m.searchResults, nil
}

func (m *MockSearchService) GetAllPostsFromIndex(ctx context.Context) ([]*domain.Post, error) {
	return nil, nil
}

func TestNewSearchHandlers(t *testing.T) {
	mockSearchService := &MockSearchService{}
	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	if handlers == nil {
		t.Fatal("Expected search handlers to be created, got nil")
	}

	if handlers.BaseHandler != base {
		t.Error("Expected base handler to be set correctly")
	}
}

func TestSearchPosts_PostMethod(t *testing.T) {
	// Create mock search service
	mockSearchService := &MockSearchService{
		searchResults: &service.SearchResponse{
			Results: []*service.SearchResult{
				{
					Post: &domain.Post{
						ID:      1,
						Title:   "Test Post",
						Image:   "test.jpg",
						Excerpt: "Test excerpt",
						Body:    "Test body",
					},
					Score: 0.95,
				},
			},
			Total:      1,
			Page:       1,
			PerPage:    10,
			TotalPages: 1,
			Query:      "test",
		},
	}

	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request
	reqBody := map[string]interface{}{
		"query":     "test",
		"page":      1,
		"per_page":  10,
		"sort_by":   "created_at:desc",
		"filter_by": "",
	}

	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPosts(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got: %d", w.Code)
	}

	// Parse response
	var response service.SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("Expected total to be 1, got: %d", response.Total)
	}

	if len(response.Results) != 1 {
		t.Errorf("Expected 1 result, got: %d", len(response.Results))
	}

	if response.Query != "test" {
		t.Errorf("Expected query to be 'test', got: %s", response.Query)
	}
}

func TestSearchPosts_WrongMethod(t *testing.T) {
	mockSearchService := &MockSearchService{}
	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request with wrong method
	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPosts(w, req)

	// Check response
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got: %d", w.Code)
	}
}

func TestSearchPosts_InvalidJSON(t *testing.T) {
	mockSearchService := &MockSearchService{}
	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPosts(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got: %d", w.Code)
	}
}

func TestSearchPosts_ServiceError(t *testing.T) {
	mockSearchService := &MockSearchService{
		searchError: domain.ErrPostNotFound,
	}

	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request
	reqBody := map[string]interface{}{
		"query": "test",
	}

	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPosts(w, req)

	// Check response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got: %d", w.Code)
	}
}

func TestSearchPostsGet_GetMethod(t *testing.T) {
	// Create mock search service
	mockSearchService := &MockSearchService{
		searchResults: &service.SearchResponse{
			Results: []*service.SearchResult{
				{
					Post: &domain.Post{
						ID:      1,
						Title:   "Test Post",
						Image:   "test.jpg",
						Excerpt: "Test excerpt",
						Body:    "Test body",
					},
					Score: 0.95,
				},
			},
			Total:      1,
			Page:       1,
			PerPage:    10,
			TotalPages: 1,
			Query:      "test",
		},
	}

	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=test&page=1&per_page=10", nil)
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPostsGet(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got: %d", w.Code)
	}

	// Parse response
	var response service.SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("Expected total to be 1, got: %d", response.Total)
	}

	if response.Query != "test" {
		t.Errorf("Expected query to be 'test', got: %s", response.Query)
	}
}

func TestSearchPostsGet_WrongMethod(t *testing.T) {
	mockSearchService := &MockSearchService{}
	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request with wrong method
	req := httptest.NewRequest(http.MethodPost, "/api/search", nil)
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPostsGet(w, req)

	// Check response
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got: %d", w.Code)
	}
}

func TestSearchPostsGet_MissingQuery(t *testing.T) {
	mockSearchService := &MockSearchService{}
	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request without query parameter
	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPostsGet(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got: %d", w.Code)
	}
}

func TestSearchPostsGet_DefaultPagination(t *testing.T) {
	mockSearchService := &MockSearchService{
		searchResults: &service.SearchResponse{
			Results:    []*service.SearchResult{},
			Total:      0,
			Page:       1,
			PerPage:    10,
			TotalPages: 0,
			Query:      "test",
		},
	}

	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request without pagination parameters
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=test", nil)
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPostsGet(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got: %d", w.Code)
	}

	// Parse response
	var response service.SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should use default values
	if response.Page != 1 {
		t.Errorf("Expected page to default to 1, got: %d", response.Page)
	}

	if response.PerPage != 10 {
		t.Errorf("Expected per_page to default to 10, got: %d", response.PerPage)
	}
}

func TestSearchPostsGet_ServiceError(t *testing.T) {
	mockSearchService := &MockSearchService{
		searchError: domain.ErrPostNotFound,
	}

	base := &BaseHandler{
		SearchService: mockSearchService,
	}

	handlers := NewSearchHandlers(base)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=test", nil)
	w := httptest.NewRecorder()

	// Call handler
	handlers.SearchPostsGet(w, req)

	// Check response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got: %d", w.Code)
	}
}
