package domain

import (
	"fmt"
	"testing"
	"time"
)

func TestNewSearchDocument(t *testing.T) {
	post := &Post{
		ID:        1,
		Title:     "Test Post",
		Image:     "test.jpg",
		Excerpt:   "Test excerpt",
		Body:      "Test body content",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	doc := NewSearchDocument(post)

	if doc.ID != fmt.Sprintf("%d", post.ID) {
		t.Errorf("Expected ID '%s', got '%s'", fmt.Sprintf("%d", post.ID), doc.ID)
	}

	if doc.Title != post.Title {
		t.Errorf("Expected Title '%s', got '%s'", post.Title, doc.Title)
	}

	if doc.Image != post.Image {
		t.Errorf("Expected Image '%s', got '%s'", post.Image, doc.Image)
	}

	if doc.Excerpt != post.Excerpt {
		t.Errorf("Expected Excerpt '%s', got '%s'", post.Excerpt, doc.Excerpt)
	}

	if doc.Body != post.Body {
		t.Errorf("Expected Body '%s', got '%s'", post.Body, doc.Body)
	}

	expectedCreatedAt := post.CreatedAt.Unix()
	if doc.CreatedAt != expectedCreatedAt {
		t.Errorf("Expected CreatedAt %d, got %d", expectedCreatedAt, doc.CreatedAt)
	}

	expectedUpdatedAt := post.UpdatedAt.Unix()
	if doc.UpdatedAt != expectedUpdatedAt {
		t.Errorf("Expected UpdatedAt %d, got %d", expectedUpdatedAt, doc.UpdatedAt)
	}
}

func TestNewSearchDocumentFromMap(t *testing.T) {
	tests := []struct {
		name        string
		data        map[string]interface{}
		expectError bool
		expectedID  int
	}{
		{
			name: "valid data with int id",
			data: map[string]interface{}{
				"id":         1,
				"title":      "Test Post",
				"body":       "Test body",
				"image":      "test.jpg",
				"excerpt":    "Test excerpt",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z",
			},
			expectError: false,
			expectedID:  1,
		},
		{
			name: "valid data with float64 id",
			data: map[string]interface{}{
				"id":         2.0,
				"title":      "Test Post 2",
				"body":       "Test body 2",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z",
			},
			expectError: false,
			expectedID:  2,
		},
		{
			name: "valid data with int64 id",
			data: map[string]interface{}{
				"id":         int64(3),
				"title":      "Test Post 3",
				"body":       "Test body 3",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z",
			},
			expectError: false,
			expectedID:  3,
		},
		{
			name: "missing id",
			data: map[string]interface{}{
				"title": "Test Post",
				"body":  "Test body",
			},
			expectError: true,
			expectedID:  0,
		},
		{
			name: "missing title",
			data: map[string]interface{}{
				"id":   1,
				"body": "Test body",
			},
			expectError: false,
			expectedID:  1,
		},
		{
			name: "missing body",
			data: map[string]interface{}{
				"id":    1,
				"title": "Test Post",
			},
			expectError: false,
			expectedID:  1,
		},
		{
			name: "invalid id type",
			data: map[string]interface{}{
				"id":    "invalid",
				"title": "Test Post",
				"body":  "Test body",
			},
			expectError: true,
			expectedID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewSearchDocumentFromMap(tt.data)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if doc.ID != fmt.Sprintf("%d", tt.expectedID) {
				t.Errorf("Expected ID '%s', got '%s'", fmt.Sprintf("%d", tt.expectedID), doc.ID)
			}

			if expectedTitle, ok := tt.data["title"]; ok {
				if doc.Title != expectedTitle {
					t.Errorf("Expected Title '%v', got '%s'", expectedTitle, doc.Title)
				}
			}

			if expectedBody, ok := tt.data["body"]; ok {
				if doc.Body != expectedBody {
					t.Errorf("Expected Body '%v', got '%s'", expectedBody, doc.Body)
				}
			}

			if image, ok := tt.data["image"]; ok {
				if doc.Image != image {
					t.Errorf("Expected Image '%v', got '%s'", image, doc.Image)
				}
			}

			if excerpt, ok := tt.data["excerpt"]; ok {
				if doc.Excerpt != excerpt {
					t.Errorf("Expected Excerpt '%v', got '%s'", excerpt, doc.Excerpt)
				}
			}

			if _, ok := tt.data["created_at"].(string); ok {
				if expectedTime, err := time.Parse(time.RFC3339, tt.data["created_at"].(string)); err == nil {
					expectedCreatedAt := expectedTime.Unix()
					if doc.CreatedAt != expectedCreatedAt {
						t.Errorf("Expected CreatedAt %d, got %d", expectedCreatedAt, doc.CreatedAt)
					}
				}
			}

			if _, ok := tt.data["updated_at"].(string); ok {
				if expectedTime, err := time.Parse(time.RFC3339, tt.data["updated_at"].(string)); err == nil {
					expectedUpdatedAt := expectedTime.Unix()
					if doc.UpdatedAt != expectedUpdatedAt {
						t.Errorf("Expected UpdatedAt %d, got %d", expectedUpdatedAt, doc.UpdatedAt)
					}
				}
			}
		})
	}
}
