package service

import (
	"context"
	"fmt"
	"testing"

	"blog-cdc-search/domain"
)

// MockSearchIndexRepositoryForSearch is a mock implementation for testing search functionality
type MockSearchIndexRepositoryForSearch struct {
	searchResults []interface{}
	searchError   error
}

func (m *MockSearchIndexRepositoryForSearch) Connect() error {
	return nil
}

func (m *MockSearchIndexRepositoryForSearch) Close() error {
	return nil
}

func (m *MockSearchIndexRepositoryForSearch) CreateCollection(schema map[string]interface{}) error {
	return nil
}

func (m *MockSearchIndexRepositoryForSearch) UpsertDocument(collectionName string, document interface{}) error {
	return nil
}

func (m *MockSearchIndexRepositoryForSearch) DeleteDocument(collectionName string, documentID string) error {
	return nil
}

func (m *MockSearchIndexRepositoryForSearch) SearchDocuments(collectionName, query string, searchParams map[string]interface{}) ([]interface{}, error) {
	if m.searchError != nil {
		return nil, m.searchError
	}
	return m.searchResults, nil
}

func (m *MockSearchIndexRepositoryForSearch) GetAllDocuments(collectionName string) ([]interface{}, error) {
	return nil, nil
}

func TestNewSearchService(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	if service == nil {
		t.Fatal("Expected search service to be created, got nil")
	}

	if service.searchRepo != mockRepo {
		t.Error("Expected search repository to be set correctly")
	}
}

func TestSearchPosts_EmptyQuery(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "",
		Page:    1,
		PerPage: 10,
	}

	result, err := service.SearchPosts(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error for empty query, got: %v", err)
	}

	if result.Total != 0 {
		t.Errorf("Expected total to be 0, got: %d", result.Total)
	}

	if len(result.Results) != 0 {
		t.Errorf("Expected no results, got: %d", len(result.Results))
	}
}

func TestSearchPosts_WhitespaceQuery(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "   ",
		Page:    1,
		PerPage: 10,
	}

	result, err := service.SearchPosts(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error for whitespace query, got: %v", err)
	}

	if result.Total != 0 {
		t.Errorf("Expected total to be 0, got: %d", result.Total)
	}
}

func TestSearchPosts_DefaultValues(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "test",
		Page:    0, // Should default to 1
		PerPage: 0, // Should default to 10
	}

	// Mock search results
	mockResults := []interface{}{
		map[string]interface{}{
			"id":          "1",
			"title":       "Test Post",
			"image":       "test.jpg",
			"excerpt":     "Test excerpt",
			"body":        "Test body",
			"created_at":  "2023-01-01 00:00:00",
			"updated_at":  "2023-01-01 00:00:00",
			"_text_match": float64(0.95),
		},
	}
	mockRepo.searchResults = mockResults

	result, err := service.SearchPosts(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Page != 1 {
		t.Errorf("Expected page to default to 1, got: %d", result.Page)
	}

	if result.PerPage != 10 {
		t.Errorf("Expected per_page to default to 10, got: %d", result.PerPage)
	}
}

func TestSearchPosts_MaxPerPage(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "test",
		Page:    1,
		PerPage: 150, // Should be capped at 100
	}

	mockRepo.searchResults = []interface{}{}

	result, err := service.SearchPosts(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.PerPage != 100 {
		t.Errorf("Expected per_page to be capped at 100, got: %d", result.PerPage)
	}
}

func TestSearchPosts_SuccessfulSearch(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "test query",
		Page:    1,
		PerPage: 10,
		SortBy:  "created_at:desc",
	}

	// Mock search results
	mockResults := []interface{}{
		map[string]interface{}{
			"id":          "1",
			"title":       "Test Post 1",
			"image":       "test1.jpg",
			"excerpt":     "Test excerpt 1",
			"body":        "Test body 1",
			"created_at":  "2023-01-01 00:00:00",
			"updated_at":  "2023-01-01 00:00:00",
			"_text_match": float64(0.95),
		},
		map[string]interface{}{
			"id":          "2",
			"title":       "Test Post 2",
			"image":       "test2.jpg",
			"excerpt":     "Test excerpt 2",
			"body":        "Test body 2",
			"created_at":  "2023-01-02 00:00:00",
			"updated_at":  "2023-01-02 00:00:00",
			"_text_match": float64(0.85),
		},
	}
	mockRepo.searchResults = mockResults

	result, err := service.SearchPosts(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected total to be 2, got: %d", result.Total)
	}

	if len(result.Results) != 2 {
		t.Errorf("Expected 2 results, got: %d", len(result.Results))
	}

	if result.Query != "test query" {
		t.Errorf("Expected query to be 'test query', got: %s", result.Query)
	}

	if result.Page != 1 {
		t.Errorf("Expected page to be 1, got: %d", result.Page)
	}

	if result.PerPage != 10 {
		t.Errorf("Expected per_page to be 10, got: %d", result.PerPage)
	}

	// Check first result
	firstResult := result.Results[0]
	if firstResult.Post.Title != "Test Post 1" {
		t.Errorf("Expected first result title to be 'Test Post 1', got: %s", firstResult.Post.Title)
	}

	if firstResult.Score != 0.95 {
		t.Errorf("Expected first result score to be 0.95, got: %f", firstResult.Score)
	}
}

func TestSearchPosts_Pagination(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "test",
		Page:    2,
		PerPage: 5,
	}

	// Mock 12 search results
	mockResults := make([]interface{}, 12)
	for i := 0; i < 12; i++ {
		mockResults[i] = map[string]interface{}{
			"id":          fmt.Sprintf("%d", i+1),
			"title":       "Test Post",
			"image":       "test.jpg",
			"excerpt":     "Test excerpt",
			"body":        "Test body",
			"created_at":  "2023-01-01 00:00:00",
			"updated_at":  "2023-01-01 00:00:00",
			"_text_match": float64(0.9),
		}
	}
	mockRepo.searchResults = mockResults

	result, err := service.SearchPosts(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedTotalPages := 3 // 12 results / 5 per page = 3 pages
	if result.TotalPages != expectedTotalPages {
		t.Errorf("Expected total pages to be %d, got: %d", expectedTotalPages, result.TotalPages)
	}
}

func TestSearchPosts_InvalidResultData(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "test",
		Page:    1,
		PerPage: 10,
	}

	// Mock invalid search results
	mockResults := []interface{}{
		"invalid result", // Not a map
		map[string]interface{}{
			"title": "Valid Post", // Missing required fields
		},
		map[string]interface{}{
			"id":          "invalid id", // Wrong type
			"title":       "Test Post",
			"image":       "test.jpg",
			"excerpt":     "Test excerpt",
			"body":        "Test body",
			"created_at":  "2023-01-01 00:00:00",
			"updated_at":  "2023-01-01 00:00:00",
			"_text_match": float64(0.95),
		},
	}
	mockRepo.searchResults = mockResults

	result, err := service.SearchPosts(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should skip invalid results and only return valid ones
	if len(result.Results) != 0 {
		t.Errorf("Expected 0 valid results (all should be invalid), got: %d", len(result.Results))
	}
}

func TestSearchPosts_RepositoryError(t *testing.T) {
	mockRepo := &MockSearchIndexRepositoryForSearch{}
	service := NewSearchService(mockRepo)

	params := SearchParams{
		Query:   "test",
		Page:    1,
		PerPage: 10,
	}

	// Mock repository error
	mockRepo.searchError = domain.ErrPostNotFound

	result, err := service.SearchPosts(context.Background(), params)
	if err == nil {
		t.Fatal("Expected error from repository, got nil")
	}

	if result != nil {
		t.Error("Expected result to be nil when error occurs")
	}
}

func TestExtractPostFromSearchResult(t *testing.T) {
	service := &SearchService{}

	// Test valid result data
	resultMap := map[string]interface{}{
		"id":         "1",
		"title":      "Test Post",
		"image":      "test.jpg",
		"excerpt":    "Test excerpt",
		"body":       "Test body",
		"created_at": "2023-01-01 00:00:00",
		"updated_at": "2023-01-01 00:00:00",
	}

	post, err := service.extractPostFromSearchResult(resultMap)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if post.ID != 1 {
		t.Errorf("Expected post ID to be 1, got: %d", post.ID)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected post title to be 'Test Post', got: %s", post.Title)
	}

	if post.Image != "test.jpg" {
		t.Errorf("Expected post image to be 'test.jpg', got: %s", post.Image)
	}

	if post.Excerpt != "Test excerpt" {
		t.Errorf("Expected post excerpt to be 'Test excerpt', got: %s", post.Excerpt)
	}

	if post.Body != "Test body" {
		t.Errorf("Expected post body to be 'Test body', got: %s", post.Body)
	}
}

func TestExtractPostFromSearchResult_MissingID(t *testing.T) {
	service := &SearchService{}

	resultMap := map[string]interface{}{
		"title": "Test Post",
		"body":  "Test body",
	}

	post, err := service.extractPostFromSearchResult(resultMap)
	if err == nil {
		t.Fatal("Expected error for missing ID, got nil")
	}

	if post != nil {
		t.Error("Expected post to be nil when error occurs")
	}
}

func TestExtractPostFromSearchResult_InvalidIDType(t *testing.T) {
	service := &SearchService{}

	resultMap := map[string]interface{}{
		"id":    "invalid",
		"title": "Test Post",
		"body":  "Test body",
	}

	post, err := service.extractPostFromSearchResult(resultMap)
	if err == nil {
		t.Fatal("Expected error for invalid ID type, got nil")
	}

	if post != nil {
		t.Error("Expected post to be nil when error occurs")
	}
}
