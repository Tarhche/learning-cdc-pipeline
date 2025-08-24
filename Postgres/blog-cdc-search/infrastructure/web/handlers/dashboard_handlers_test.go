package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDashboardHandlers_ServeDashboard(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewDashboardHandlers(base)

	// Create test posts
	mockService.CreatePost(context.Background(), "Admin Post 1", "image1.jpg", "excerpt1", "body1")
	mockService.CreatePost(context.Background(), "Admin Post 2", "image2.jpg", "excerpt2", "body2")

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	recorder := httptest.NewRecorder()

	handler.ServeDashboard(recorder, req)

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
	if !strings.Contains(body, "Blog Admin Dashboard") {
		t.Error("response should contain dashboard title")
	}

	if !strings.Contains(body, "Admin Post 1") {
		t.Error("response should contain first post title")
	}

	if !strings.Contains(body, "Admin Post 2") {
		t.Error("response should contain second post title")
	}

	if !strings.Contains(body, "Create New Post") {
		t.Error("response should contain create new post link")
	}

	if !strings.Contains(body, "Edit") {
		t.Error("response should contain edit buttons")
	}

	if !strings.Contains(body, "Delete") {
		t.Error("response should contain delete buttons")
	}

	if !strings.Contains(body, "Back to Blog") {
		t.Error("response should contain back to blog link")
	}
}

func TestDashboardHandlers_ServeDashboard_EmptyPosts(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewDashboardHandlers(base)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	recorder := httptest.NewRecorder()

	handler.ServeDashboard(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "No posts yet") {
		t.Error("response should contain 'No posts yet' message")
	}

	if !strings.Contains(body, "Create your first blog post") {
		t.Error("response should contain create first post message")
	}
}

func TestDashboardHandlers_ServeCreateForm(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewDashboardHandlers(base)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/create", nil)
	recorder := httptest.NewRecorder()

	handler.ServeCreateForm(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	// Check content type
	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected HTML content type, got %s", contentType)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "Create New Post") {
		t.Error("response should contain create form title")
	}

	if !strings.Contains(body, `<form id="createForm">`) {
		t.Error("response should contain create form")
	}

	if !strings.Contains(body, `name="title"`) {
		t.Error("response should contain title input")
	}

	if !strings.Contains(body, `name="image"`) {
		t.Error("response should contain image input")
	}

	if !strings.Contains(body, `name="excerpt"`) {
		t.Error("response should contain excerpt input")
	}

	if !strings.Contains(body, `name="body"`) {
		t.Error("response should contain body input")
	}

	if !strings.Contains(body, "Back to Dashboard") {
		t.Error("response should contain back to dashboard link")
	}
}

func TestDashboardHandlers_ServeEditForm(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewDashboardHandlers(base)

	// Create a test post
	post, _ := mockService.CreatePost(context.Background(), "Edit Test Post", "edit_image.jpg", "Edit excerpt", "Edit body content")

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
			url := "/dashboard/edit"
			if tt.postID != "" {
				url += "?id=" + tt.postID
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			handler.ServeEditForm(recorder, req)

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
				if !strings.Contains(body, "Edit Post") {
					t.Error("response should contain edit form title")
				}

				if !strings.Contains(body, `<form id="editForm">`) {
					t.Error("response should contain edit form")
				}

				if !strings.Contains(body, post.Title) {
					t.Error("response should contain pre-filled post title")
				}

				if !strings.Contains(body, post.Image) {
					t.Error("response should contain pre-filled post image")
				}

				if !strings.Contains(body, post.Excerpt) {
					t.Error("response should contain pre-filled post excerpt")
				}

				if !strings.Contains(body, post.Body) {
					t.Error("response should contain pre-filled post body")
				}

				if !strings.Contains(body, "Update Post") {
					t.Error("response should contain update button")
				}

				if !strings.Contains(body, "Back to Dashboard") {
					t.Error("response should contain back to dashboard link")
				}
			}
		})
	}
}

func TestDashboardHandlers_MethodNotAllowed(t *testing.T) {
	mockService := NewMockPostService()
	base := &BaseHandler{PostService: mockService}
	handler := NewDashboardHandlers(base)

	tests := []struct {
		name    string
		method  string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:    "dashboard POST",
			method:  http.MethodPost,
			handler: handler.ServeDashboard,
		},
		{
			name:    "create form POST",
			method:  http.MethodPost,
			handler: handler.ServeCreateForm,
		},
		{
			name:    "edit form POST",
			method:  http.MethodPost,
			handler: handler.ServeEditForm,
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
