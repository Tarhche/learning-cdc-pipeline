package domain

import (
	"fmt"
	"time"
)

type SearchDocument struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Excerpt   string `json:"excerpt"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func NewSearchDocument(post *Post) *SearchDocument {
	return &SearchDocument{
		ID:        fmt.Sprintf("%d", post.ID),
		Title:     post.Title,
		Image:     post.Image,
		Excerpt:   post.Excerpt,
		Body:      post.Body,
		CreatedAt: post.CreatedAt.Unix(),
		UpdatedAt: post.UpdatedAt.Unix(),
	}
}

func NewSearchDocumentFromMap(data map[string]interface{}) (*SearchDocument, error) {
	id, ok := data["id"]
	if !ok {
		return nil, ErrPostNotFound
	}

	var postID int
	switch v := id.(type) {
	case float64:
		postID = int(v)
	case int:
		postID = v
	case int64:
		postID = int(v)
	default:
		return nil, ErrPostNotFound
	}

	title, _ := data["title"].(string)
	image, _ := data["image"].(string)
	excerpt, _ := data["excerpt"].(string)
	body, _ := data["body"].(string)

	var createdAt, updatedAt int64

	if createdVal, ok := data["created_at"]; ok {
		switch v := createdVal.(type) {
		case float64:
			createdAt = int64(v)
		case int64:
			createdAt = v
		case int:
			createdAt = int64(v)
		case string:
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				createdAt = t.Unix()
			}
		}
	}

	if updatedVal, ok := data["updated_at"]; ok {
		switch v := updatedVal.(type) {
		case float64:
			updatedAt = int64(v)
		case int64:
			updatedAt = v
		case int:
			updatedAt = int64(v)
		case string:
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				updatedAt = t.Unix()
			}
		}
	}

	if createdAt == 0 {
		createdAt = time.Now().Unix()
	}
	if updatedAt == 0 {
		updatedAt = time.Now().Unix()
	}

	return &SearchDocument{
		ID:        fmt.Sprintf("%d", postID),
		Title:     title,
		Image:     image,
		Excerpt:   excerpt,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
