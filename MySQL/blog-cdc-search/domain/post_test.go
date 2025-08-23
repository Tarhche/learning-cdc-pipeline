package domain

import (
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		image   string
		excerpt string
		body    string
		wantErr bool
	}{
		{
			name:    "Valid post",
			title:   "Test Title",
			image:   "https://example.com/image.jpg",
			excerpt: "Test excerpt",
			body:    "Test body content",
			wantErr: false,
		},
		{
			name:    "Empty title",
			title:   "",
			image:   "https://example.com/image.jpg",
			excerpt: "Test excerpt",
			body:    "Test body content",
			wantErr: true,
		},
		{
			name:    "Empty body",
			title:   "Test Title",
			image:   "https://example.com/image.jpg",
			excerpt: "Test excerpt",
			body:    "",
			wantErr: true,
		},
		{
			name:    "Empty title and body",
			title:   "",
			image:   "https://example.com/image.jpg",
			excerpt: "Test excerpt",
			body:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPost(tt.title, tt.image, tt.excerpt, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewPost() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("NewPost() unexpected error: %v", err)
				return
			}

			if post.Title != tt.title {
				t.Errorf("NewPost() title = %v, want %v", post.Title, tt.title)
			}

			if post.Image != tt.image {
				t.Errorf("NewPost() image = %v, want %v", post.Image, tt.image)
			}

			if post.Excerpt != tt.excerpt {
				t.Errorf("NewPost() excerpt = %v, want %v", post.Excerpt, tt.excerpt)
			}

			if post.Body != tt.body {
				t.Errorf("NewPost() body = %v, want %v", post.Body, tt.body)
			}

			if post.ID != 0 {
				t.Errorf("NewPost() ID = %v, want 0", post.ID)
			}

			// Check that timestamps are set
			if post.CreatedAt.IsZero() {
				t.Error("NewPost() CreatedAt should not be zero")
			}

			if post.UpdatedAt.IsZero() {
				t.Error("NewPost() UpdatedAt should not be zero")
			}
		})
	}
}

func TestPost_Update(t *testing.T) {
	post, err := NewPost("Original Title", "original.jpg", "Original excerpt", "Original body")
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	originalCreatedAt := post.CreatedAt
	originalUpdatedAt := post.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	err = post.Update("New Title", "new.jpg", "New excerpt", "New body")
	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}

	if post.Title != "New Title" {
		t.Errorf("Update() title = %v, want 'New Title'", post.Title)
	}

	if post.Image != "new.jpg" {
		t.Errorf("Update() image = %v, want 'new.jpg'", post.Image)
	}

	if post.Excerpt != "New excerpt" {
		t.Errorf("Update() excerpt = %v, want 'New excerpt'", post.Excerpt)
	}

	if post.Body != "New body" {
		t.Errorf("Update() body = %v, want 'New body'", post.Body)
	}

	if post.CreatedAt != originalCreatedAt {
		t.Errorf("Update() should not change CreatedAt")
	}

	if post.UpdatedAt == originalUpdatedAt {
		t.Errorf("Update() should update UpdatedAt")
	}
}

func TestPost_Update_Validation(t *testing.T) {
	post, err := NewPost("Original Title", "original.jpg", "Original excerpt", "Original body")
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		title   string
		image   string
		excerpt string
		body    string
		wantErr bool
	}{
		{
			name:    "Empty title",
			title:   "",
			image:   "new.jpg",
			excerpt: "New excerpt",
			body:    "New body",
			wantErr: true,
		},
		{
			name:    "Empty body",
			title:   "New Title",
			image:   "new.jpg",
			excerpt: "New excerpt",
			body:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := post.Update(tt.title, tt.image, tt.excerpt, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Update() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Update() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPost_Validate(t *testing.T) {
	tests := []struct {
		name    string
		post    *Post
		wantErr bool
	}{
		{
			name: "Valid post",
			post: &Post{
				Title: "Test Title",
				Body:  "Test body",
			},
			wantErr: false,
		},
		{
			name: "Empty title",
			post: &Post{
				Title: "",
				Body:  "Test body",
			},
			wantErr: true,
		},
		{
			name: "Empty body",
			post: &Post{
				Title: "Test Title",
				Body:  "",
			},
			wantErr: true,
		},
		{
			name: "Empty title and body",
			post: &Post{
				Title: "",
				Body:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}
