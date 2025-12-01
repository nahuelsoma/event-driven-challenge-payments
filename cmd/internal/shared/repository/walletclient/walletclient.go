package walletclient

import (
	"context"
	"errors"
	"log/slog"

	"net/http"
)

// WalletClient implements all wallet operations
type WalletClient struct {
	client *http.Client
}

// NewWalletClient creates a new WalletClient
func NewWalletClient(client *http.Client) (*WalletClient, error) {
	if client == nil {
		return nil, errors.New("wallet client: client cannot be nil")
	}

	return &WalletClient{client: client}, nil
}

// Reserve reserves funds in the wallet for a payment
func (wc *WalletClient) Reserve(ctx context.Context, userID string, amount float64, paymentID string) error {
	slog.DebugContext(ctx, "[DEBUG] WalletClient.Reserve called", "user_id", userID, "amount", amount, "payment_id", paymentID)
	// TODO: Implement the logic to reserve the funds via HTTP client
	// POST /api/v1/wallets/:user_id/reserve
	return nil
}

// Confirm confirms the reserved funds deduction in the wallet
func (wc *WalletClient) Confirm(ctx context.Context, userID string, amount float64, paymentID string) error {
	slog.DebugContext(ctx, "[DEBUG] WalletClient.Confirm called", "user_id", userID, "amount", amount, "payment_id", paymentID)
	// TODO: Implement the logic to confirm the funds via HTTP client
	// POST /api/v1/wallets/:user_id/confirm
	return nil
}

// Release releases reserved funds back to available balance
func (wc *WalletClient) Release(ctx context.Context, userID string, amount float64, paymentID string) error {
	slog.DebugContext(ctx, "[DEBUG] WalletClient.Release called", "user_id", userID, "amount", amount, "payment_id", paymentID)
	// TODO: Implement the logic to release the funds via HTTP client
	// POST /api/v1/wallets/:user_id/release
	return nil
}
