package creator

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// Errors
var (
	ErrNilWalletClient = errors.New("wallet client cannot be nil")
)

// PaymentRequest represents a request to create a payment
type PaymentRequest struct {
	UserID   string          `json:"user_id"`  // User ID of the payment (in a real application, this would be the user ID from the authenticated user)
	Amount   float64         `json:"amount"`   // Amount of the payment
	Currency domain.Currency `json:"currency"` // Currency of the payment
}

// Validate validates the payment request
// It returns an error if the request is invalid
func (p *PaymentRequest) Validate() error {
	if p.UserID == "" {
		return errors.New("user ID is required")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if err := p.Currency.Validate(); err != nil {
		return err
	}
	return nil
}

// NewPayment creates a new payment
func NewPayment(idempotencyKey string, userID string, amount float64, currency domain.Currency) *domain.Payment {
	return &domain.Payment{
		ID:             uuid.New().String(),
		IdempotencyKey: idempotencyKey,
		UserID:         userID,
		Amount:         amount,
		Currency:       currency,
		Status:         domain.StatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
