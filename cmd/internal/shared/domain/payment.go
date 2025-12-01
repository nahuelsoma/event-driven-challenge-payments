package domain

import (
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
