package repository

import (
	"context"
	"database/sql"
	"fmt"

	"blog-cdc-search/domain"
)

// DBExecutor defines the interface for database operations
type DBExecutor interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// PostgreSQLPostRepository implements PostRepository using PostgreSQL
type PostgreSQLPostRepository struct {
	db DBExecutor
}

// NewPostgreSQLPostRepository creates a new PostgreSQLPostRepository instance
func NewPostgreSQLPostRepository(db DBExecutor) *PostgreSQLPostRepository {
	return &PostgreSQLPostRepository{db: db}
}

// Create inserts a new post into the database
func (r *PostgreSQLPostRepository) Create(ctx context.Context, post *domain.Post) error {
	query := `
		INSERT INTO posts (title, image, excerpt, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query, post.Title, post.Image, post.Excerpt, post.Body, post.CreatedAt, post.UpdatedAt).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	post.ID = id
	return nil
}

// GetByID retrieves a post by its ID
func (r *PostgreSQLPostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	query := `
		SELECT id, title, image, excerpt, body, created_at, updated_at
		FROM posts WHERE id = $1
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
func (r *PostgreSQLPostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
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
func (r *PostgreSQLPostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `
		UPDATE posts 
		SET title = $1, image = $2, excerpt = $3, body = $4, updated_at = $5
		WHERE id = $6
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
func (r *PostgreSQLPostRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = $1`

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
