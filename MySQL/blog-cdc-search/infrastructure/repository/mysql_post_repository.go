package repository

import (
	"context"
	"database/sql"
	"fmt"

	"blog-cdc-search/domain"
)

// MySQLPostRepository implements PostRepository using MySQL
type MySQLPostRepository struct {
	db *sql.DB
}

// NewMySQLPostRepository creates a new MySQLPostRepository instance
func NewMySQLPostRepository(db *sql.DB) *MySQLPostRepository {
	return &MySQLPostRepository{db: db}
}

// Create inserts a new post into the database
func (r *MySQLPostRepository) Create(ctx context.Context, post *domain.Post) error {
	query := `
		INSERT INTO posts (title, image, excerpt, body, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, post.Title, post.Image, post.Excerpt, post.Body, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	post.ID = int(id)
	return nil
}

// GetByID retrieves a post by its ID
func (r *MySQLPostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	query := `
		SELECT id, title, image, excerpt, body, created_at, updated_at
		FROM posts WHERE id = ?
	`

	var post domain.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.Title, &post.Image, &post.Excerpt, &post.Body, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// GetAll retrieves all posts from the database
func (r *MySQLPostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	query := `
		SELECT id, title, image, excerpt, body, created_at, updated_at
		FROM posts ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Image, &post.Excerpt, &post.Body, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return posts, nil
}

// Update updates an existing post in the database
func (r *MySQLPostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `
		UPDATE posts 
		SET title = ?, image = ?, excerpt = ?, body = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, post.Title, post.Image, post.Excerpt, post.Body, post.UpdatedAt, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID: %d", post.ID)
	}

	return nil
}

// Delete removes a post from the database
func (r *MySQLPostRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID: %d", id)
	}

	return nil
}
