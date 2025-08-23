package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"blog-cdc-search/domain"

	"github.com/gorilla/mux"
)

// MockSearchServiceForWeb implements SearchService interface for testing
type MockSearchServiceForWeb struct {
	posts []*domain.Post
}

func NewMockSearchServiceForWeb() *MockSearchServiceForWeb {
	return &MockSearchServiceForWeb{
		posts: make([]*domain.Post, 0),
	}
}

func (m *MockSearchServiceForWeb) SearchPosts(ctx context.Context, params interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockSearchServiceForWeb) GetAllPostsFromIndex(ctx context.Context) ([]*domain.Post, error) {
	return m.posts, nil
}

func (m *MockSearchServiceForWeb) AddPost(post *domain.Post) {
	m.posts = append(m.posts, post)
}

func TestWebHandlers_ServeHomePage(t *testing.T) {
	mockPostService := NewMockPostService()
	mockSearchService := NewMockSearchServiceForWeb()

	base := &BaseHandler{
		PostService:   mockPostService,
		SearchService: mockSearchService,
	}
	handler := NewWebHandlers(base)

	// Create test posts
	post1, _ := mockPostService.CreatePost(context.Background(), "Post 1", "image1.jpg", "excerpt1", "body1")
	post2, _ := mockPostService.CreatePost(context.Background(), "Post 2", "image2.jpg", "excerpt2", "body2")

	// Add posts to search service
	mockSearchService.AddPost(post1)
	mockSearchService.AddPost(post2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHomePage(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	// Check content type
	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected HTML content type, got %s", contentType)
	}

	// Check that the response contains expected content
	body := recorder.Body.String()
	if !strings.Contains(body, "My Blog") {
		t.Error("response should contain blog title")
	}

	if !strings.Contains(body, "Post 1") {
		t.Error("response should contain first post title")
	}

	if !strings.Contains(body, "Post 2") {
		t.Error("response should contain second post title")
	}

	if !strings.Contains(body, "Admin Dashboard") {
		t.Error("response should contain admin dashboard link")
	}
}

func TestWebHandlers_ServeHomePage_EmptyPosts(t *testing.T) {
	mockPostService := NewMockPostService()
	mockSearchService := NewMockSearchServiceForWeb()

	base := &BaseHandler{
		PostService:   mockPostService,
		SearchService: mockSearchService,
	}
	handler := NewWebHandlers(base)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHomePage(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "No posts yet") {
		t.Error("response should contain 'No posts yet' message")
	}
}

func TestWebHandlers_ServePostDetail(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewWebHandlers(base)

	// Create a test post
	mockService.CreatePost(context.Background(), "Test Post", "image.jpg", "Test excerpt", "Test body content")

	tests := []struct {
		name           string
		postID         string
		expectedStatus int
		expectedTitle  string
	}{
		{
			name:           "existing post",
			postID:         "1",
			expectedStatus: http.StatusOK,
			expectedTitle:  "Test Post",
		},
		{
			name:           "non-existing post",
			postID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedTitle:  "",
		},
		{
			name:           "invalid post ID",
			postID:         "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/post/"+tt.postID, nil)

			// Set up mux vars to simulate gorilla/mux path variables
			req = mux.SetURLVars(req, map[string]string{"id": tt.postID})

			recorder := httptest.NewRecorder()

			handler.ServePostDetail(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				// Check content type
				contentType := recorder.Header().Get("Content-Type")
				if !strings.Contains(contentType, "text/html") {
					t.Errorf("expected HTML content type, got %s", contentType)
				}

				body := recorder.Body.String()
				if !strings.Contains(body, tt.expectedTitle) {
					t.Errorf("response should contain post title '%s'", tt.expectedTitle)
				}

				if !strings.Contains(body, "Test body content") {
					t.Error("response should contain post body")
				}

				if !strings.Contains(body, "Back to Blog") {
					t.Error("response should contain back link")
				}
			}
		})
	}
}

func TestWebHandlers_MethodNotAllowed(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewWebHandlers(base)

	tests := []struct {
		name    string
		method  string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:    "home page POST",
			method:  http.MethodPost,
			handler: handler.ServeHomePage,
		},
		{
			name:    "post detail POST",
			method:  http.MethodPost,
			handler: handler.ServePostDetail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			recorder := httptest.NewRecorder()

			tt.handler(recorder, req)

			if recorder.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
			}
		})
	}
}
