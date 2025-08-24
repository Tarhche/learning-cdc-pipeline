package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"blog-cdc-search/domain"
)

// Mock implementations for testing
type MockMessageQueueRepository struct {
	connectCalled bool
	closeCalled   bool
	consumeCalled bool
	handler       func([]byte) error
	connectError  error
}

func (m *MockMessageQueueRepository) Connect() error {
	m.connectCalled = true
	return m.connectError
}

func (m *MockMessageQueueRepository) Close() error {
	m.closeCalled = true
	return nil
}

func (m *MockMessageQueueRepository) ConsumeMessages(queueName string, handler func([]byte) error) error {
	m.consumeCalled = true
	m.handler = handler
	return nil
}

func (m *MockMessageQueueRepository) PublishMessage(exchange, routingKey string, message []byte) error {
	return nil
}

type MockSearchIndexRepository struct {
	connectCalled          bool
	closeCalled            bool
	createCollectionCalled bool
	upsertDocumentCalled   bool
	deleteDocumentCalled   bool
	connectError           error
	createCollectionError  error
}

func (m *MockSearchIndexRepository) Connect() error {
	m.connectCalled = true
	return m.connectError
}

func (m *MockSearchIndexRepository) Close() error {
	m.closeCalled = true
	return nil
}

func (m *MockSearchIndexRepository) CreateCollection(schema map[string]interface{}) error {
	m.createCollectionCalled = true
	return m.createCollectionError
}

func (m *MockSearchIndexRepository) UpsertDocument(collectionName string, document interface{}) error {
	m.upsertDocumentCalled = true
	return nil
}

func (m *MockSearchIndexRepository) DeleteDocument(collectionName string, documentID string) error {
	m.deleteDocumentCalled = true
	return nil
}

func (m *MockSearchIndexRepository) SearchDocuments(collectionName, query string, searchParams map[string]interface{}) ([]interface{}, error) {
	return nil, nil
}

func (m *MockSearchIndexRepository) GetAllDocuments(collectionName string) ([]interface{}, error) {
	return nil, nil
}

func TestNewCDCService(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if service.messageQueue != mockMQ {
		t.Error("Expected messageQueue to be set")
	}

	if service.searchIndex != mockSearch {
		t.Error("Expected searchIndex to be set")
	}
}

func TestConvertTimestamps(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}
	service := NewCDCService(mockMQ, mockSearch)

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "nanosecond timestamps",
			input: map[string]interface{}{
				"id":         1,
				"created_at": 1756051794247827, // Nanoseconds timestamp
				"updated_at": 1756051794247827, // Nanoseconds timestamp
			},
			expected: map[string]interface{}{
				"id":         1,
				"created_at": int64(1756051794247827 / 1000000000), // Converted to seconds
				"updated_at": int64(1756051794247827 / 1000000000), // Converted to seconds
			},
		},
		{
			name: "normal timestamps",
			input: map[string]interface{}{
				"id":         2,
				"created_at": 1756057014, // Normal Unix timestamp
				"updated_at": 1756057014, // Normal Unix timestamp
			},
			expected: map[string]interface{}{
				"id":         2,
				"created_at": 1756057014, // Should remain unchanged
				"updated_at": 1756057014, // Should remain unchanged
			},
		},
		{
			name: "mixed timestamps",
			input: map[string]interface{}{
				"id":         3,
				"created_at": 1756051794247827, // Nanoseconds timestamp
				"updated_at": 1756057014,       // Normal Unix timestamp
			},
			expected: map[string]interface{}{
				"id":         3,
				"created_at": int64(1756051794247827 / 1000000000), // Converted to seconds
				"updated_at": 1756057014,                           // Should remain unchanged
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the input to avoid modifying the original
			inputCopy := make(map[string]interface{})
			for k, v := range tt.input {
				inputCopy[k] = v
			}

			// Call the convertTimestamps method
			service.convertTimestamps(inputCopy)

			// Check created_at
			if createdAt, ok := inputCopy["created_at"]; ok {
				expectedCreatedAt := tt.expected["created_at"]
				if createdAt != expectedCreatedAt {
					t.Errorf("created_at: expected %v, got %v", expectedCreatedAt, createdAt)
				}
			}

			// Check updated_at
			if updatedAt, ok := inputCopy["updated_at"]; ok {
				expectedUpdatedAt := tt.expected["updated_at"]
				if updatedAt != expectedUpdatedAt {
					t.Errorf("updated_at: expected %v, got %v", expectedUpdatedAt, updatedAt)
				}
			}
		})
	}
}

func TestCDCService_StartCDC_Success(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)
	ctx, cancel := context.WithCancel(context.Background())

	// Start the service in a goroutine since it blocks
	go func() {
		err := service.StartCDC(ctx, "test-queue")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	}()

	// Give it a moment to start and then cancel
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Wait a bit for cleanup
	time.Sleep(50 * time.Millisecond)

	if !mockMQ.connectCalled {
		t.Error("Expected Connect to be called on message queue")
	}

	if !mockSearch.connectCalled {
		t.Error("Expected Connect to be called on search index")
	}

	if !mockSearch.createCollectionCalled {
		t.Error("Expected CreateCollection to be called")
	}

	if !mockMQ.consumeCalled {
		t.Error("Expected ConsumeMessages to be called")
	}
}

func TestCDCService_StartCDC_MessageQueueConnectError(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{
		connectError: errors.New("connection failed"),
	}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := service.StartCDC(ctx, "test-queue")
	if err == nil {
		t.Error("Expected error when message queue connection fails")
	}

	if err.Error() != "failed to connect to RabbitMQ: connection failed" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCDCService_StartCDC_SearchIndexConnectError(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{
		connectError: errors.New("connection failed"),
	}

	service := NewCDCService(mockMQ, mockSearch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := service.StartCDC(ctx, "test-queue")
	if err == nil {
		t.Error("Expected error when search index connection fails")
	}

	if err.Error() != "failed to connect to Typesense: connection failed" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCDCService_StartCDC_CreateCollectionError(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{
		createCollectionError: errors.New("create collection failed"),
	}

	service := NewCDCService(mockMQ, mockSearch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the service in a goroutine since it blocks
	go func() {
		err := service.StartCDC(ctx, "test-queue")
		if err != nil {
			t.Errorf("Expected no error when collection creation fails (continues gracefully), got: %v", err)
		}
	}()

	// Give it a moment to start and then cancel
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Wait a bit for cleanup
	time.Sleep(50 * time.Millisecond)

	if !mockMQ.connectCalled {
		t.Error("Expected Connect to be called on message queue")
	}

	if !mockSearch.connectCalled {
		t.Error("Expected Connect to be called on search index")
	}

	if !mockSearch.createCollectionCalled {
		t.Error("Expected CreateCollection to be called")
	}

	// The service should continue even when collection creation fails
	if !mockMQ.consumeCalled {
		t.Error("Expected ConsumeMessages to be called even after collection creation failure")
	}
}

func TestCDCService_handleMessage_ValidInsertEvent(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	// Create a valid insert event
	event := &domain.CDCEvent{
		Database: "blog",
		Table:    "posts",
		Type:     domain.EventTypeInsert,
		Data: map[string]interface{}{
			"id":         1,
			"title":      "Test Post",
			"body":       "Test body",
			"created_at": "2023-01-01 00:00:00",
			"updated_at": "2023-01-01 00:00:00",
		},
	}

	// Convert to JSON
	message, err := event.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert event to JSON: %v", err)
	}

	// Handle the message
	err = service.handleMessage(message)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !mockSearch.upsertDocumentCalled {
		t.Error("Expected UpsertDocument to be called")
	}
}

func TestCDCService_handleMessage_ValidUpdateEvent(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	// Create a valid update event
	event := &domain.CDCEvent{
		Database: "blog",
		Table:    "posts",
		Type:     domain.EventTypeUpdate,
		Data: map[string]interface{}{
			"id":         1,
			"title":      "Updated Post",
			"body":       "Updated body",
			"created_at": "2023-01-01 00:00:00",
			"updated_at": "2023-01-01 00:00:00",
		},
	}

	// Convert to JSON
	message, err := event.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert event to JSON: %v", err)
	}

	// Handle the message
	err = service.handleMessage(message)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !mockSearch.upsertDocumentCalled {
		t.Error("Expected UpsertDocument to be called")
	}
}

func TestCDCService_handleMessage_ValidDeleteEvent(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	// Create a valid delete event
	event := &domain.CDCEvent{
		Database: "blog",
		Table:    "posts",
		Type:     domain.EventTypeDelete,
		Data: map[string]interface{}{
			"id": 1,
		},
	}

	// Convert to JSON
	message, err := event.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert event to JSON: %v", err)
	}

	// Handle the message
	err = service.handleMessage(message)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !mockSearch.deleteDocumentCalled {
		t.Error("Expected DeleteDocument to be called")
	}
}

func TestCDCService_handleMessage_InvalidJSON(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	// Invalid JSON message
	invalidMessage := []byte(`{"invalid": json}`)

	err := service.handleMessage(invalidMessage)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestCDCService_handleMessage_InvalidEvent(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	// Invalid event (missing required fields)
	event := &domain.CDCEvent{
		Database: "blog",
		// Missing table and type
		Data: map[string]interface{}{"id": 1},
	}

	message, err := event.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert event to JSON: %v", err)
	}

	err = service.handleMessage(message)
	if err == nil {
		t.Error("Expected error for invalid event")
	}
}

func TestCDCService_handleMessage_NonPostsTable(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	// Event for a different table
	event := &domain.CDCEvent{
		Database: "blog",
		Table:    "users", // Different table
		Type:     domain.EventTypeInsert,
		Data:     map[string]interface{}{"id": 1},
	}

	message, err := event.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert event to JSON: %v", err)
	}

	// Should not error, just skip
	err = service.handleMessage(message)
	if err != nil {
		t.Errorf("Expected no error for non-posts table, got: %v", err)
	}

	// Should not call any search index methods
	if mockSearch.upsertDocumentCalled || mockSearch.deleteDocumentCalled {
		t.Error("Expected no search index operations for non-posts table")
	}
}

func TestCDCService_handleMessage_UnknownEventType(t *testing.T) {
	mockMQ := &MockMessageQueueRepository{}
	mockSearch := &MockSearchIndexRepository{}

	service := NewCDCService(mockMQ, mockSearch)

	// Unknown event type
	event := &domain.CDCEvent{
		Database: "blog",
		Table:    "posts",
		Type:     "unknown", // Unknown type
		Data:     map[string]interface{}{"id": 1},
	}

	message, err := event.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert event to JSON: %v", err)
	}

	err = service.handleMessage(message)
	if err == nil {
		t.Error("Expected error for unknown event type")
	}

	if err.Error() != "unknown event type: unknown" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}
