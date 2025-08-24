package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"blog-cdc-search/domain"
)

// SearchService handles search operations
type SearchService struct {
	searchRepo domain.SearchIndexRepository
}

// SearchResult represents a search result with relevance score
type SearchResult struct {
	Post       *domain.Post        `json:"post"`
	Score      float64             `json:"score"`
	Highlights map[string][]string `json:"highlights,omitempty"`
}

// SearchParams represents search parameters
type SearchParams struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	SortBy   string `json:"sort_by"`
	FilterBy string `json:"filter_by"`
}

// SearchResponse represents the complete search response
type SearchResponse struct {
	Results    []*SearchResult `json:"results"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
	Query      string          `json:"query"`
}

// NewSearchService creates a new search service instance
func NewSearchService(searchRepo domain.SearchIndexRepository) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
	}
}

// SearchPosts performs a search for posts based on the given parameters
func (s *SearchService) SearchPosts(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	// Validate search parameters
	if strings.TrimSpace(params.Query) == "" {
		return &SearchResponse{
			Results:    []*SearchResult{},
			Total:      0,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: 0,
			Query:      params.Query,
		}, nil
	}

	// Validate and set default values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 10
	}
	if params.PerPage > 100 {
		params.PerPage = 100
	}

	// Prepare search parameters for Typesense
	searchParams := map[string]interface{}{
		"page":     params.Page,
		"per_page": params.PerPage,
		"sort_by":  "_text_match:desc,created_at:desc",
		"query_by": "title,excerpt,body",
	}

	if params.FilterBy != "" {
		searchParams["filter_by"] = params.FilterBy
	}

	// Perform search
	results, err := s.searchRepo.SearchDocuments("posts", params.Query, searchParams)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert results to SearchResult format
	searchResults := make([]*SearchResult, 0, len(results))
	for _, result := range results {
		// Convert the interface{} result to a map
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract post data
		post, err := s.extractPostFromSearchResult(resultMap)
		if err != nil {
			continue // Skip invalid results
		}

		// Extract score
		score := 0.0
		if scoreVal, ok := resultMap["_text_match"].(float64); ok {
			score = scoreVal
		}

		// Extract highlights if available
		highlights := make(map[string][]string)
		if highlightsVal, ok := resultMap["highlights"].(map[string]interface{}); ok {
			for field, highlightList := range highlightsVal {
				if highlightArray, ok := highlightList.([]interface{}); ok {
					highlights[field] = make([]string, 0, len(highlightArray))
					for _, highlight := range highlightArray {
						if highlightStr, ok := highlight.(string); ok {
							highlights[field] = append(highlights[field], highlightStr)
						}
					}
				}
			}
		}

		searchResults = append(searchResults, &SearchResult{
			Post:       post,
			Score:      score,
			Highlights: highlights,
		})
	}

	// Calculate pagination info
	total := len(searchResults)
	totalPages := (total + params.PerPage - 1) / params.PerPage

	return &SearchResponse{
		Results:    searchResults,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
		Query:      params.Query,
	}, nil
}

// GetAllPostsFromIndex retrieves all posts from the search index
func (s *SearchService) GetAllPostsFromIndex(ctx context.Context) ([]*domain.Post, error) {
	// Get all documents from the posts collection
	results, err := s.searchRepo.GetAllDocuments("posts")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts from index: %w", err)
	}

	// Convert results to Post domain objects
	posts := make([]*domain.Post, 0, len(results))
	for _, result := range results {
		// Convert the interface{} result to a map
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract post data
		post, err := s.extractPostFromSearchResult(resultMap)
		if err != nil {
			continue // Skip invalid results
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// extractPostFromSearchResult converts a search result map to a Post domain object
func (s *SearchService) extractPostFromSearchResult(resultMap map[string]interface{}) (*domain.Post, error) {
	// Extract ID
	idVal, ok := resultMap["id"]
	if !ok {
		return nil, fmt.Errorf("missing post ID")
	}

	var id int
	switch v := idVal.(type) {
	case string:
		if parsedID, err := strconv.Atoi(v); err == nil {
			id = parsedID
		} else {
			return nil, fmt.Errorf("invalid post ID format: %s", v)
		}
	case float64:
		id = int(v)
	case int:
		id = v
	case int64:
		id = int(v)
	default:
		return nil, fmt.Errorf("invalid post ID type")
	}

	// Extract other fields
	title, _ := resultMap["title"].(string)
	image, _ := resultMap["image"].(string)
	excerpt, _ := resultMap["excerpt"].(string)
	body, _ := resultMap["body"].(string)

	// Extract and convert timestamps from Unix timestamps to time.Time
	var createdAt, updatedAt time.Time

	if createdVal, ok := resultMap["created_at"]; ok {
		switch v := createdVal.(type) {
		case float64:
			createdAt = time.Unix(int64(v), 0)
		case int64:
			createdAt = time.Unix(v, 0)
		case int:
			createdAt = time.Unix(int64(v), 0)
		}
	}

	if updatedVal, ok := resultMap["updated_at"]; ok {
		switch v := updatedVal.(type) {
		case float64:
			updatedAt = time.Unix(int64(v), 0)
		case int64:
			updatedAt = time.Unix(v, 0)
		case int:
			updatedAt = time.Unix(int64(v), 0)
		}
	}

	// Create post object
	post := &domain.Post{
		ID:        id,
		Title:     title,
		Image:     image,
		Excerpt:   excerpt,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return post, nil
}
