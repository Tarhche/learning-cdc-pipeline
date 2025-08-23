package handlers

import (
	"blog-cdc-search/domain"
	"context"
)

// PostService interface for mocking in tests
type PostService interface {
	CreatePost(ctx context.Context, title, image, excerpt, body string) (*domain.Post, error)
	GetPost(ctx context.Context, id int) (*domain.Post, error)
	GetAllPosts(ctx context.Context) ([]*domain.Post, error)
	UpdatePost(ctx context.Context, id int, title, image, excerpt, body string) (*domain.Post, error)
	DeletePost(ctx context.Context, id int) error
}

// SearchService interface for mocking in tests
type SearchService interface {
	SearchPosts(ctx context.Context, params interface{}) (interface{}, error)
	GetAllPostsFromIndex(ctx context.Context) ([]*domain.Post, error)
}

// BaseHandler contains common dependencies
type BaseHandler struct {
	PostService   PostService
	SearchService SearchService
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(postService PostService, searchService SearchService) *BaseHandler {
	return &BaseHandler{
		PostService:   postService,
		SearchService: searchService,
	}
}
