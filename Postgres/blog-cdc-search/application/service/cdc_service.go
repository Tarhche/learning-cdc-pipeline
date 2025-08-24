package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"blog-cdc-search/domain"
)

// CDCService handles the business logic for CDC events
type CDCService struct {
	messageQueue domain.MessageQueueRepository
	searchIndex  domain.SearchIndexRepository
}

// NewCDCService creates a new CDC service instance
func NewCDCService(messageQueue domain.MessageQueueRepository, searchIndex domain.SearchIndexRepository) *CDCService {
	return &CDCService{
		messageQueue: messageQueue,
		searchIndex:  searchIndex,
	}
}

// StartCDC starts the CDC pipeline
func (s *CDCService) StartCDC(ctx context.Context, queueName string) error {
	log.Printf("Starting CDC service, listening to queue: %s", queueName)

	// Connect to RabbitMQ
	if err := s.messageQueue.Connect(); err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer s.messageQueue.Close()

	// Connect to Typesense√
	if err := s.searchIndex.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Typesense: %w", err)
	}
	defer s.searchIndex.Close()

	// Ensure the posts collection exists
	if err := s.ensurePostsCollection(); err != nil {
		return fmt.Errorf("failed to ensure posts collection: %w", err)
	}

	// Start consuming messages
	return s.messageQueue.ConsumeMessages(queueName, s.handleMessage)
}

// handleMessage processes individual messages from RabbitMQ
func (s *CDCService) handleMessage(message []byte) error {
	// Parse the CDC event (try Debezium format first, fallback to original format)
	event, err := domain.FromDebeziumJSON(message)
	if err != nil {
		log.Printf("Failed to parse CDC event: %v", err)
		return err
	}

	// Validate the event
	if !event.IsValid() {
		log.Printf("Invalid CDC event: %+v", event)
		return fmt.Errorf("invalid CDC event")
	}

	// Only process posts table events
	if event.Table != "posts" {
		log.Printf("Skipping non-posts table event: %s.%s", event.Database, event.Table)
		return nil
	}

	// Process the event based on its type
	switch event.Type {
	case domain.EventTypeInsert, domain.EventTypeUpdate, domain.EventTypeBootstrapInsert:
		return s.handleUpsert(event)
	case domain.EventTypeDelete:
		return s.handleDelete(event)
	case domain.EventTypeBootstrapStart:
		log.Printf("Bootstrap insert event received")
		return nil
	case domain.EventTypeBootstrapComplete:
		log.Printf("Bootstrap complete event received")
		return nil
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// convertTimestamps converts nanosecond timestamps to Unix timestamps in the event data
func (s *CDCService) convertTimestamps(data map[string]interface{}) {
	// Convert created_at if present
	if createdVal, ok := data["created_at"]; ok {
		switch v := createdVal.(type) {
		case float64:
			if v > 1e12 { // If greater than 1 trillion, likely nanoseconds
				// Use time package to convert nanoseconds to Unix timestamp
				t := time.Unix(0, int64(v))
				data["created_at"] = t.Unix() // Extract Unix timestamp in seconds
			}
		case int64:
			if v > 1e12 { // If greater than 1 trillion, likely nanoseconds
				// Use time package to convert nanoseconds to Unix timestamp
				t := time.Unix(0, v)
				data["created_at"] = t.Unix() // Extract Unix timestamp in seconds
			}
		case int:
			if v > 1e12 { // If greater than 1 trillion, likely nanoseconds
				// Use time package to convert nanoseconds to Unix timestamp
				t := time.Unix(0, int64(v))
				data["created_at"] = t.Unix() // Extract Unix timestamp in seconds
			}
		}
	}

	// Convert updated_at if present
	if updatedVal, ok := data["updated_at"]; ok {
		switch v := updatedVal.(type) {
		case float64:
			if v > 1e12 { // If greater than 1 trillion, likely nanoseconds
				// Use time package to convert nanoseconds to Unix timestamp
				t := time.Unix(0, int64(v))
				data["updated_at"] = t.Unix() // Extract Unix timestamp in seconds
			}
		case int64:
			if v > 1e12 { // If greater than 1 trillion, likely nanoseconds
				// Use time package to convert nanoseconds to Unix timestamp
				t := time.Unix(0, v)
				data["updated_at"] = t.Unix() // Extract Unix timestamp in seconds
			}
		case int:
			if v > 1e12 { // If greater than 1 trillion, likely nanoseconds
				// Use time package to convert nanoseconds to Unix timestamp
				t := time.Unix(0, int64(v))
				data["updated_at"] = t.Unix() // Extract Unix timestamp in seconds
			}
		}
	}
}

// handleUpsert processes insert, update, and bootstrap events
func (s *CDCService) handleUpsert(event *domain.CDCEvent) error {
	eventType := event.Type
	log.Printf("Processing %s event for post ID: %v", eventType, event.Data["id"])

	// Convert timestamps from nanoseconds to Unix timestamps before processing
	s.convertTimestamps(event.Data)

	// Create search document from the event data
	doc, err := domain.NewSearchDocumentFromMap(event.Data)
	if err != nil {
		log.Printf("Failed to create search document: %v", err)
		return err
	}

	// Upsert the document to Typesense
	collectionName := "posts"
	if err := s.searchIndex.UpsertDocument(collectionName, doc); err != nil {
		log.Printf("Failed to upsert document to Typesense: %v", err)
		return err
	}

	log.Printf("Successfully indexed post ID: %s", doc.ID)
	return nil
}

// handleDelete processes delete events
func (s *CDCService) handleDelete(event *domain.CDCEvent) error {
	log.Printf("Processing delete event for post ID: %v", event.Data["id"])

	// Get the ID from the event data
	id, exists := event.GetID()
	if !exists {
		log.Printf("Failed to extract ID from delete event")
		return fmt.Errorf("failed to extract ID from delete event")
	}

	// Delete the document from Typesense
	collectionName := "posts"
	documentID := fmt.Sprintf("%d", id)
	if err := s.searchIndex.DeleteDocument(collectionName, documentID); err != nil {
		log.Printf("Failed to delete document from Typesense: %v", err)
		return err
	}

	log.Printf("Successfully removed post ID: %d from index", id)
	return nil
}

// ensurePostsCollection ensures the posts collection exists in Typesense
func (s *CDCService) ensurePostsCollection() error {
	schema := map[string]interface{}{
		"name": "posts",
		"fields": []map[string]interface{}{
			{"name": "id", "type": "string"},
			{"name": "title", "type": "string", "facet": false, "index": true},
			{"name": "image", "type": "string", "optional": true, "facet": false, "index": false},
			{"name": "excerpt", "type": "string", "optional": true, "facet": false, "index": true},
			{"name": "body", "type": "string", "facet": false, "index": true},
			{"name": "created_at", "type": "int64", "facet": false, "index": true},
			{"name": "updated_at", "type": "int64", "facet": false, "index": false},
		},
		"default_sorting_field": "created_at",
	}

	// Try to create the collection first
	err := s.searchIndex.CreateCollection(schema)
	if err != nil {
		// If collection already exists, we need to ensure it has the right schema
		// For now, we'll just log this and continue
		// In a production environment, you might want to handle schema updates
		log.Printf("Collection creation failed (might already exist): %v", err)
	}

	// Note: In Typesense, when you upsert a document with the same ID field value,
	// it will replace the existing document. The id field is used as the document
	// identifier to ensure proper updates.

	return nil
}
