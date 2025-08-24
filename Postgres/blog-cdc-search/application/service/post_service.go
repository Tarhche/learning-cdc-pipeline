package service

import (
	"context"
	"errors"

	"blog-cdc-search/domain"
)

// PostService handles business logic for posts
type PostService struct {
	repo domain.PostRepository
}

// NewPostService creates a new PostService instance
func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(ctx context.Context, title, image, excerpt, body string) (*domain.Post, error) {
	post, err := domain.NewPost(title, image, excerpt, body)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// GetPost retrieves a post by ID
func (s *PostService) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	if id <= 0 {
		return nil, errors.New("invalid post ID")
	}

	return s.repo.GetByID(ctx, id)
}

// GetAllPosts retrieves all posts
func (s *PostService) GetAllPosts(ctx context.Context) ([]*domain.Post, error) {
	return s.repo.GetAll(ctx)
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(ctx context.Context, id int, title, image, excerpt, body string) (*domain.Post, error) {
	if id <= 0 {
		return nil, errors.New("invalid post ID")
	}

	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := post.Update(title, image, excerpt, body); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid post ID")
	}

	return s.repo.Delete(ctx, id)
}
