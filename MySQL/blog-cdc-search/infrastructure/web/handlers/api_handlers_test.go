package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"blog-cdc-search/domain"
)

// MockPostService implements PostService interface for testing
type MockPostService struct {
	posts  map[int]*domain.Post
	nextID int
}

func NewMockPostService() *MockPostService {
	return &MockPostService{
		posts:  make(map[int]*domain.Post),
		nextID: 1,
	}
}

func (m *MockPostService) CreatePost(ctx context.Context, title, image, excerpt, body string) (*domain.Post, error) {
	post := &domain.Post{
		ID:        m.nextID,
		Title:     title,
		Image:     image,
		Excerpt:   excerpt,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.posts[post.ID] = post
	m.nextID++
	return post, nil
}

func (m *MockPostService) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return nil, domain.ErrPostNotFound
	}
	return post, nil
}

func (m *MockPostService) GetAllPosts(ctx context.Context) ([]*domain.Post, error) {
	posts := make([]*domain.Post, 0, len(m.posts))
	for _, post := range m.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *MockPostService) UpdatePost(ctx context.Context, id int, title, image, excerpt, body string) (*domain.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return nil, domain.ErrPostNotFound
	}

	post.Title = title
	post.Image = image
	post.Excerpt = excerpt
	post.Body = body
	post.UpdatedAt = time.Now()

	return post, nil
}

func (m *MockPostService) DeletePost(ctx context.Context, id int) error {
	_, exists := m.posts[id]
	if !exists {
		return domain.ErrPostNotFound
	}
	delete(m.posts, id)
	return nil
}

func TestAPIHandlers_CreatePost(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewAPIHandlers(base)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: map[string]string{
				"title":   "Test Title",
				"image":   "http://example.com/image.jpg",
				"excerpt": "Test excerpt",
				"body":    "Test body content",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid JSON",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody []byte
			if tt.requestBody != nil {
				requestBody, _ = json.Marshal(tt.requestBody)
			} else {
				requestBody = []byte("invalid json")
			}

			req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler.CreatePost(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response domain.Post
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				if response.Title != tt.requestBody["title"] {
					t.Errorf("expected title %s, got %s", tt.requestBody["title"], response.Title)
				}
			}
		})
	}
}

func TestAPIHandlers_GetPost(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewAPIHandlers(base)

	// Create a test post
	post, _ := mockService.CreatePost(context.Background(), "Test Title", "image.jpg", "excerpt", "body")

	tests := []struct {
		name           string
		postID         string
		expectedStatus int
	}{
		{
			name:           "existing post",
			postID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existing post",
			postID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid post ID",
			postID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing post ID",
			postID:         "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/posts"
			if tt.postID != "" {
				url += "?id=" + tt.postID
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			handler.GetPost(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.Post
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				if response.ID != post.ID {
					t.Errorf("expected post ID %d, got %d", post.ID, response.ID)
				}
			}
		})
	}
}

func TestAPIHandlers_GetAllPosts(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewAPIHandlers(base)

	// Create test posts
	mockService.CreatePost(context.Background(), "Post 1", "image1.jpg", "excerpt1", "body1")
	mockService.CreatePost(context.Background(), "Post 2", "image2.jpg", "excerpt2", "body2")

	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	recorder := httptest.NewRecorder()

	handler.GetAllPosts(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response []*domain.Post
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("expected 2 posts, got %d", len(response))
	}
}

func TestAPIHandlers_UpdatePost(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewAPIHandlers(base)

	// Create a test post
	mockService.CreatePost(context.Background(), "Original Title", "image.jpg", "excerpt", "body")

	tests := []struct {
		name           string
		postID         string
		requestBody    map[string]string
		expectedStatus int
	}{
		{
			name:   "successful update",
			postID: "1",
			requestBody: map[string]string{
				"title":   "Updated Title",
				"image":   "updated_image.jpg",
				"excerpt": "Updated excerpt",
				"body":    "Updated body",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existing post",
			postID:         "999",
			requestBody:    map[string]string{"title": "Updated"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid post ID",
			postID:         "invalid",
			requestBody:    map[string]string{"title": "Updated"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tt.requestBody)
			url := "/api/posts?id=" + tt.postID

			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler.UpdatePost(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.Post
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}

				if response.Title != tt.requestBody["title"] {
					t.Errorf("expected title %s, got %s", tt.requestBody["title"], response.Title)
				}
			}
		})
	}
}

func TestAPIHandlers_DeletePost(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewAPIHandlers(base)

	// Create a test post
	mockService.CreatePost(context.Background(), "Test Title", "image.jpg", "excerpt", "body")

	tests := []struct {
		name           string
		postID         string
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			postID:         "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "non-existing post",
			postID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid post ID",
			postID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/posts?id=" + tt.postID

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			recorder := httptest.NewRecorder()

			handler.DeletePost(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}
