package domain

import (
	"encoding/json"
	"testing"
)

func TestNewCDCEvent(t *testing.T) {
	data := map[string]interface{}{
		"id":    1,
		"title": "Test Post",
		"body":  "Test Body",
	}

	event := NewCDCEvent("blog", "posts", EventTypeInsert, data)

	if event.Database != "blog" {
		t.Errorf("Expected database 'blog', got '%s'", event.Database)
	}

	if event.Table != "posts" {
		t.Errorf("Expected table 'posts', got '%s'", event.Table)
	}

	if event.Type != EventTypeInsert {
		t.Errorf("Expected type '%s', got '%s'", EventTypeInsert, event.Type)
	}

	if event.Data == nil {
		t.Error("Expected data to be set")
	}

	if event.TS == 0 {
		t.Error("Expected timestamp to be set")
	}
}

func TestCDCEvent_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		event    *CDCEvent
		expected bool
	}{
		{
			name: "valid event",
			event: &CDCEvent{
				Database: "blog",
				Table:    "posts",
				Type:     EventTypeInsert,
				Data:     map[string]interface{}{"id": 1},
			},
			expected: true,
		},
		{
			name: "missing database",
			event: &CDCEvent{
				Table: EventTypeInsert,
				Type:  EventTypeInsert,
				Data:  map[string]interface{}{"id": 1},
			},
			expected: false,
		},
		{
			name: "missing table",
			event: &CDCEvent{
				Database: "blog",
				Type:     EventTypeInsert,
				Data:     map[string]interface{}{"id": 1},
			},
			expected: false,
		},
		{
			name: "missing type",
			event: &CDCEvent{
				Database: "blog",
				Table:    "posts",
				Data:     map[string]interface{}{"id": 1},
			},
			expected: false,
		},
		{
			name: "missing data",
			event: &CDCEvent{
				Database: "blog",
				Table:    "posts",
				Type:     EventTypeInsert,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.IsValid()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCDCEvent_GetID(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected int
		exists   bool
	}{
		{
			name:     "int id",
			data:     map[string]interface{}{"id": 1},
			expected: 1,
			exists:   true,
		},
		{
			name:     "float64 id",
			data:     map[string]interface{}{"id": 2.0},
			expected: 2,
			exists:   true,
		},
		{
			name:     "int64 id",
			data:     map[string]interface{}{"id": int64(3)},
			expected: 3,
			exists:   true,
		},
		{
			name:     "missing id",
			data:     map[string]interface{}{"title": "test"},
			expected: 0,
			exists:   false,
		},
		{
			name:     "invalid id type",
			data:     map[string]interface{}{"id": "invalid"},
			expected: 0,
			exists:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &CDCEvent{Data: tt.data}
			id, exists := event.GetID()

			if exists != tt.exists {
				t.Errorf("Expected exists %v, got %v", tt.exists, exists)
			}

			if exists && id != tt.expected {
				t.Errorf("Expected id %d, got %d", tt.expected, id)
			}
		})
	}
}

func TestCDCEvent_ToJSON(t *testing.T) {
	event := &CDCEvent{
		Database: "blog",
		Table:    "posts",
		Type:     EventTypeInsert,
		Data:     map[string]interface{}{"id": 1, "title": "test"},
		TS:       1234567890,
	}

	jsonData, err := event.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	// Verify it's valid JSON by parsing it back
	var parsedEvent CDCEvent
	if err := json.Unmarshal(jsonData, &parsedEvent); err != nil {
		t.Fatalf("Failed to parse JSON back: %v", err)
	}

	if parsedEvent.Database != event.Database {
		t.Errorf("Expected database '%s', got '%s'", event.Database, parsedEvent.Database)
	}

	if parsedEvent.Table != event.Table {
		t.Errorf("Expected table '%s', got '%s'", event.Table, parsedEvent.Table)
	}

	if parsedEvent.Type != event.Type {
		t.Errorf("Expected type '%s', got '%s'", event.Type, parsedEvent.Type)
	}
}

func TestFromJSON(t *testing.T) {
	jsonData := `{
		"database": "blog",
		"table": "posts",
		"type": "insert",
		"data": {"id": 1, "title": "test"},
		"ts": 1234567890
	}`

	event, err := FromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if event.Database != "blog" {
		t.Errorf("Expected database 'blog', got '%s'", event.Database)
	}

	if event.Table != "posts" {
		t.Errorf("Expected table 'posts', got '%s'", event.Table)
	}

	if event.Type != "insert" {
		t.Errorf("Expected type 'insert', got '%s'", event.Type)
	}

	if event.Data["id"] != float64(1) {
		t.Errorf("Expected id 1, got %v", event.Data["id"])
	}

	if event.Data["title"] != "test" {
		t.Errorf("Expected title 'test', got '%v'", event.Data["title"])
	}
}

func TestFromJSON_Invalid(t *testing.T) {
	invalidJSON := `{"invalid": json}`

	_, err := FromJSON([]byte(invalidJSON))
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}
