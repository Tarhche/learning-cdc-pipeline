package service

import (
	"context"
	"errors"
	"testing"

	"blog-cdc-search/domain"
)

// MockPostRepository is a mock implementation of PostRepository for testing
type MockPostRepository struct {
	posts  map[int]*domain.Post
	nextID int
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[int]*domain.Post),
		nextID: 1,
	}
}

func (m *MockPostRepository) Create(ctx context.Context, post *domain.Post) error {
	post.ID = m.nextID
	m.posts[post.ID] = post
	m.nextID++
	return nil
}

func (m *MockPostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *MockPostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	posts := make([]*domain.Post, 0, len(m.posts))
	for _, post := range m.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *MockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	if _, exists := m.posts[post.ID]; !exists {
		return errors.New("post not found")
	}
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Delete(ctx context.Context, id int) error {
	if _, exists := m.posts[id]; !exists {
		return errors.New("post not found")
	}
	delete(m.posts, id)
	return nil
}

func TestPostService_CreatePost(t *testing.T) {
	repo := NewMockPostRepository()
	service := NewPostService(repo)
	ctx := context.Background()

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := service.CreatePost(ctx, tt.title, tt.image, tt.excerpt, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreatePost() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("CreatePost() unexpected error: %v", err)
				return
			}

			if post.Title != tt.title {
				t.Errorf("CreatePost() title = %v, want %v", post.Title, tt.title)
			}

			if post.Body != tt.body {
				t.Errorf("CreatePost() body = %v, want %v", post.Body, tt.body)
			}

			if post.ID == 0 {
				t.Error("CreatePost() should set post ID")
			}
		})
	}
}

func TestPostService_GetPost(t *testing.T) {
	repo := NewMockPostRepository()
	service := NewPostService(repo)
	ctx := context.Background()

	// Create a test post first
	post, err := service.CreatePost(ctx, "Test Title", "test.jpg", "Test excerpt", "Test body")
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Valid ID",
			id:      post.ID,
			wantErr: false,
		},
		{
			name:    "Invalid ID - zero",
			id:      0,
			wantErr: true,
		},
		{
			name:    "Invalid ID - negative",
			id:      -1,
			wantErr: true,
		},
		{
			name:    "Non-existent ID",
			id:      999,
			wantErr: false, // Repository will handle the error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetPost(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPost() expected error but got none")
				}
				return
			}

			if tt.id == 999 {
				// Non-existent ID should return repository error
				if err == nil {
					t.Errorf("GetPost() expected error for non-existent ID")
				}
				return
			}

			if err != nil {
				t.Errorf("GetPost() unexpected error: %v", err)
				return
			}

			if result.ID != tt.id {
				t.Errorf("GetPost() ID = %v, want %v", result.ID, tt.id)
			}
		})
	}
}

func TestPostService_GetAllPosts(t *testing.T) {
	repo := NewMockPostRepository()
	service := NewPostService(repo)
	ctx := context.Background()

	// Initially no posts
	posts, err := service.GetAllPosts(ctx)
	if err != nil {
		t.Errorf("GetAllPosts() unexpected error: %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("GetAllPosts() expected 0 posts, got %d", len(posts))
	}

	// Create some test posts
	_, err = service.CreatePost(ctx, "Post 1", "img1.jpg", "Excerpt 1", "Body 1")
	if err != nil {
		t.Fatalf("Failed to create test post 1: %v", err)
	}

	_, err = service.CreatePost(ctx, "Post 2", "img2.jpg", "Excerpt 2", "Body 2")
	if err != nil {
		t.Fatalf("Failed to create test post 2: %v", err)
	}

	// Now should have 2 posts
	posts, err = service.GetAllPosts(ctx)
	if err != nil {
		t.Errorf("GetAllPosts() unexpected error: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("GetAllPosts() expected 2 posts, got %d", len(posts))
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	repo := NewMockPostRepository()
	service := NewPostService(repo)
	ctx := context.Background()

	// Create a test post first
	post, err := service.CreatePost(ctx, "Original Title", "original.jpg", "Original excerpt", "Original body")
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		title   string
		image   string
		excerpt string
		body    string
		wantErr bool
	}{
		{
			name:    "Valid update",
			id:      post.ID,
			title:   "Updated Title",
			image:   "updated.jpg",
			excerpt: "Updated excerpt",
			body:    "Updated body",
			wantErr: false,
		},
		{
			name:    "Invalid ID - zero",
			id:      0,
			title:   "Updated Title",
			image:   "updated.jpg",
			excerpt: "Updated excerpt",
			body:    "Updated body",
			wantErr: true,
		},
		{
			name:    "Empty title",
			id:      post.ID,
			title:   "",
			image:   "updated.jpg",
			excerpt: "Updated excerpt",
			body:    "Updated body",
			wantErr: true,
		},
		{
			name:    "Empty body",
			id:      post.ID,
			title:   "Updated Title",
			image:   "updated.jpg",
			excerpt: "Updated excerpt",
			body:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.UpdatePost(ctx, tt.id, tt.title, tt.image, tt.excerpt, tt.body)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdatePost() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdatePost() unexpected error: %v", err)
				return
			}

			if result.Title != tt.title {
				t.Errorf("UpdatePost() title = %v, want %v", result.Title, tt.title)
			}

			if result.Body != tt.body {
				t.Errorf("UpdatePost() body = %v, want %v", result.Body, tt.body)
			}
		})
	}
}

func TestPostService_DeletePost(t *testing.T) {
	repo := NewMockPostRepository()
	service := NewPostService(repo)
	ctx := context.Background()

	// Create a test post first
	post, err := service.CreatePost(ctx, "Test Title", "test.jpg", "Test excerpt", "Test body")
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Valid ID",
			id:      post.ID,
			wantErr: false,
		},
		{
			name:    "Invalid ID - zero",
			id:      0,
			wantErr: true,
		},
		{
			name:    "Invalid ID - negative",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeletePost(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeletePost() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DeletePost() unexpected error: %v", err)
				return
			}

			// Verify post was deleted
			if tt.id > 0 {
				_, err := service.GetPost(ctx, tt.id)
				if err == nil {
					t.Error("DeletePost() should have deleted the post")
				}
			}
		})
	}
}
