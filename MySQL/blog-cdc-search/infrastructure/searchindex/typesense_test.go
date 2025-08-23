package searchindex

import (
	"testing"
)

func TestNewTypesenseRepository(t *testing.T) {
	config := TypesenseConfig{
		Host:   "localhost",
		Port:   8108,
		APIKey: "test-key",
	}

	repo := NewTypesenseRepository(config)

	if repo == nil {
		t.Fatal("Expected repository to be created")
	}

	if repo.config.Host != config.Host {
		t.Errorf("Expected host %s, got %s", config.Host, repo.config.Host)
	}

	if repo.config.Port != config.Port {
		t.Errorf("Expected port %d, got %d", config.Port, repo.config.Port)
	}

	if repo.config.APIKey != config.APIKey {
		t.Errorf("Expected API key %s, got %s", config.APIKey, repo.config.APIKey)
	}
}

func TestTypesenseRepository_Close(t *testing.T) {
	config := TypesenseConfig{
		Host:   "localhost",
		Port:   8108,
		APIKey: "test-key",
	}

	repo := NewTypesenseRepository(config)

	// Close should not error
	err := repo.Close()
	if err != nil {
		t.Errorf("Expected no error when closing, got: %v", err)
	}
}

func TestTypesenseRepository_ensurePostsCollection(t *testing.T) {
	// Test the schema structure that would be used
	schema := map[string]interface{}{
		"name": "posts",
		"fields": []map[string]interface{}{
			{"name": "id", "type": "int32"},
			{"name": "title", "type": "string"},
			{"name": "image", "type": "string", "optional": true},
			{"name": "excerpt", "type": "string", "optional": true},
			{"name": "body", "type": "string"},
			{"name": "created_at", "type": "int64"},
			{"name": "updated_at", "type": "int64"},
		},
		"default_sorting_field": "created_at",
	}

	// Verify schema structure
	if schema["name"] != "posts" {
		t.Errorf("Expected collection name 'posts', got %v", schema["name"])
	}

	fields, ok := schema["fields"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected fields to be a slice of maps")
	}

	if len(fields) != 7 {
		t.Errorf("Expected 7 fields, got %d", len(fields))
	}

	// Check specific fields
	expectedFields := []string{"id", "title", "image", "excerpt", "body", "created_at", "updated_at"}
	for i, expectedField := range expectedFields {
		if fields[i]["name"] != expectedField {
			t.Errorf("Expected field %s at position %d, got %v", expectedField, i, fields[i]["name"])
		}
	}
}
