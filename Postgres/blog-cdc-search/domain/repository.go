package domain

import "context"

// PostRepository defines the interface for post data access
type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int) (*Post, error)
	GetAll(ctx context.Context) ([]*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error
}

// MessageQueueRepository defines the interface for RabbitMQ operations
type MessageQueueRepository interface {
	Connect() error
	Close() error
	ConsumeMessages(queueName string, handler func([]byte) error) error
	PublishMessage(exchange, routingKey string, message []byte) error
}

// SearchIndexRepository defines the interface for Typesense operations
type SearchIndexRepository interface {
	Connect() error
	Close() error
	CreateCollection(schema map[string]interface{}) error
	UpsertDocument(collectionName string, document interface{}) error
	DeleteDocument(collectionName string, documentID string) error
	SearchDocuments(collectionName, query string, searchParams map[string]interface{}) ([]interface{}, error)
	GetAllDocuments(collectionName string) ([]interface{}, error)
}
