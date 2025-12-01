package creator

import (
	"context"
	"errors"
)

type WalletReserverRepository struct {
	client interface{}
}

func NewWalletReserverRepository(client interface{}) (*WalletReserverRepository, error) {
	if client == nil {
		return nil, errors.New("payment storer: client cannot be nil")
	}

	return &WalletReserverRepository{client: client}, nil
}

func (r *WalletReserverRepository) Reserve(ctx context.Context, userID string, amount float64, paymentID string) error {
	// TODO: Implement the logic to reserve the funds
	return nil
}
