package domain

import (
	"encoding/json"
	"errors"
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID             string    `json:"id"`              // Unique identifier for the payment
	IdempotencyKey string    `json:"idempotency_key"` // Idempotency key for the payment
	UserID         string    `json:"user_id"`         // User ID of the payment owner
	Amount         float64   `json:"amount"`          // Amount of the payment
	Currency       Currency  `json:"currency"`        // Currency of the payment
	Status         Status    `json:"status"`          // Status of the payment
	CreatedAt      time.Time `json:"created_at"`      // Timestamp when the payment was created
	UpdatedAt      time.Time `json:"updated_at"`      // Timestamp when the payment was updated
}

// Validate validates the payment
// It validates the payment and returns an error if the validation fails
func (p *Payment) Validate() error {
	if p.ID == "" {
		return errors.New("payment ID is required")
	}
	if p.UserID == "" {
		return errors.New("user ID is required")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

// UpdateStatus updates the status of the payment
// It updates the status of the payment and the updated at timestamp
func (p *Payment) UpdateStatus(status Status) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

// Parse parses a payment from bytes
// It parses a payment from bytes and returns an error if the parsing fails
func (p *Payment) Parse(body []byte) error {
	return json.Unmarshal(body, p)
}

// Marshal marshals a payment to bytes
// It marshals a payment to bytes and returns an error if the marshalling fails
func (p *Payment) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
