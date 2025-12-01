package processor

import (
	"context"
	"errors"
	"log/slog"
)

// WalletConfirmerRepository confirms funds in the wallet service
type WalletConfirmerRepository struct {
	client interface{}
}

// NewWalletConfirmerRepository creates a new WalletConfirmerRepository
func NewWalletConfirmerRepository(client interface{}) (*WalletConfirmerRepository, error) {
	if client == nil {
		return nil, errors.New("wallet confirmer: client cannot be nil")
	}

	return &WalletConfirmerRepository{client: client}, nil
}

// Confirm confirms funds in the wallet for a payment
func (r *WalletConfirmerRepository) Confirm(ctx context.Context, userID string, amount float64, paymentID string) error {
	slog.DebugContext(ctx, "[DEBUG] WalletConfirmerRepository.Confirm called", "user_id", userID, "amount", amount, "payment_id", paymentID)
	// TODO: Implement the logic to confirm the funds via HTTP client
	return nil
}
