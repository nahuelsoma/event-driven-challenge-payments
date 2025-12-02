package walletclient

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockWalletClient is a mock implementation of WalletClient for testing
type MockWalletClient struct {
	mock.Mock
}

// Reserve reserves funds in the wallet for a payment
func (m *MockWalletClient) Reserve(ctx context.Context, userID string, amount float64, paymentID string) error {
	args := m.Called(ctx, userID, amount, paymentID)
	return args.Error(0)
}

// Confirm confirms the reserved funds deduction in the wallet
func (m *MockWalletClient) Confirm(ctx context.Context, userID string, amount float64, paymentID string) error {
	args := m.Called(ctx, userID, amount, paymentID)
	return args.Error(0)
}

// Release releases reserved funds back to available balance
func (m *MockWalletClient) Release(ctx context.Context, userID string, amount float64, paymentID string) error {
	args := m.Called(ctx, userID, amount, paymentID)
	return args.Error(0)
}
