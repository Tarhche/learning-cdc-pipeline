package domain

import (
	"encoding/json"
	"time"
)

// CDCEvent represents a Change Data Capture event from Maxwell
type CDCEvent struct {
	Database string                 `json:"database"`
	Table    string                 `json:"table"`
	Type     string                 `json:"type"` // insert, update, delete
	Data     map[string]interface{} `json:"data"`
	Old      map[string]interface{} `json:"old,omitempty"`
	TS       int64                  `json:"ts"`
	Xid      int64                  `json:"xid,omitempty"`
	Xoffset  int64                  `json:"xoffset,omitempty"`
}

// EventType constants
const (
	EventTypeInsert            = "insert"
	EventTypeUpdate            = "update"
	EventTypeDelete            = "delete"
	EventTypeBootstrapStart    = "bootstrap-start"
	EventTypeBootstrapInsert   = "bootstrap-insert"
	EventTypeBootstrapComplete = "bootstrap-complete"
)

// NewCDCEvent creates a new CDC event
func NewCDCEvent(database, table, eventType string, data map[string]interface{}) *CDCEvent {
	return &CDCEvent{
		Database: database,
		Table:    table,
		Type:     eventType,
		Data:     data,
		TS:       time.Now().Unix(),
	}
}

// IsValid checks if the CDC event is valid
func (e *CDCEvent) IsValid() bool {
	return e.Database != "" && e.Table != "" && e.Type != "" && e.Data != nil
}

// GetID extracts the ID from the data
func (e *CDCEvent) GetID() (int, bool) {
	if id, exists := e.Data["id"]; exists {
		switch v := id.(type) {
		case float64:
			return int(v), true
		case int:
			return v, true
		case int64:
			return int(v), true
		}
	}
	return 0, false
}

// ToJSON converts the event to JSON
func (e *CDCEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON creates a CDC event from JSON
func FromJSON(data []byte) (*CDCEvent, error) {
	var event CDCEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
