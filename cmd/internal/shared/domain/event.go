package domain

import (
	"encoding/json"
	"time"
)

// Event represents a payment event in the event store
type Event struct {
	ID        string          `json:"id"`         // Unique identifier for the event
	PaymentID string          `json:"payment_id"` // Payment ID associated with the event
	Sequence  int             `json:"sequence"`   // Sequence number of the event for this payment
	EventType string          `json:"event_type"` // Type of the event (created, reserved, completed, failed)
	Payload   json.RawMessage `json:"payload"`    // Event payload as JSON
	CreatedAt time.Time       `json:"created_at"` // Timestamp when the event was created
}
