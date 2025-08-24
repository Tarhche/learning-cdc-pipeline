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

// FromDebeziumJSON creates a CDC event from Debezium JSON format
func FromDebeziumJSON(data []byte) (*CDCEvent, error) {
	// Try to parse as Debezium format first
	var debeziumEvent struct {
		Payload struct {
			Before map[string]interface{} `json:"before"`
			After  map[string]interface{} `json:"after"`
			Source struct {
				DB    string `json:"db"`
				Table string `json:"table"`
			} `json:"source"`
			Op string `json:"op"`
			TS int64  `json:"ts_ms"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(data, &debeziumEvent); err != nil {
		// If it's not Debezium format, try the original format
		return FromJSON(data)
	}

	// Debug: Check if we actually have valid Debezium data
	if debeziumEvent.Payload.Source.DB == "" && debeziumEvent.Payload.Source.Table == "" {
		// This doesn't look like valid Debezium data, fall back to original format
		return FromJSON(data)
	}

	// Convert Debezium format to CDCEvent format
	event := &CDCEvent{
		Database: debeziumEvent.Payload.Source.DB,
		Table:    debeziumEvent.Payload.Source.Table,
		TS:       debeziumEvent.Payload.TS,
	}

	// Map Debezium operation types to our event types
	switch debeziumEvent.Payload.Op {
	case "r": // read (snapshot)
		event.Type = EventTypeBootstrapInsert
		event.Data = debeziumEvent.Payload.After
	case "c": // create
		event.Type = EventTypeInsert
		event.Data = debeziumEvent.Payload.After
	case "u": // update
		event.Type = EventTypeUpdate
		event.Data = debeziumEvent.Payload.After
		event.Old = debeziumEvent.Payload.Before
	case "d": // delete
		event.Type = EventTypeDelete
		event.Data = debeziumEvent.Payload.Before // Use 'before' for delete
	default:
		event.Type = debeziumEvent.Payload.Op
		event.Data = debeziumEvent.Payload.After
	}

	return event, nil
}
