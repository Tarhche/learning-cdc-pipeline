package domain

import (
	"errors"
	"time"
)

// Domain errors
var (
	ErrPostNotFound = errors.New("post not found")
)

// Post represents a blog post entity
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Excerpt   string    `json:"excerpt"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPost creates a new Post instance
func NewPost(title, image, excerpt, body string) (*Post, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if body == "" {
		return nil, errors.New("body cannot be empty")
	}

	now := time.Now()
	return &Post{
		Title:     title,
		Image:     image,
		Excerpt:   excerpt,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update updates the post fields
func (p *Post) Update(title, image, excerpt, body string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	if body == "" {
		return errors.New("body cannot be empty")
	}

	p.Title = title
	p.Image = image
	p.Excerpt = excerpt
	p.Body = body
	p.UpdatedAt = time.Now()
	return nil
}

// Validate checks if the post is valid
func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("title cannot be empty")
	}
	if p.Body == "" {
		return errors.New("body cannot be empty")
	}
	return nil
}
