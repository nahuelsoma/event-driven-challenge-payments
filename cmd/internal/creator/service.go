package creator

import (
	"context"
	"errors"
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// WalletReserver interface for reserving funds
type WalletReserver interface {
	Reserve(ctx context.Context, userID string, amount float64, paymentID string) error
}

// PaymentStorer interface for storing payments
type PaymentStorer interface {
	GetByIDempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error)
	Save(ctx context.Context, payment *domain.Payment) error
	UpdateStatus(ctx context.Context, paymentID string, status domain.Status, gatewayRef string) error
}

// PaymentPublisher interface for publishing payments
type PaymentPublisher interface {
	Publish(ctx context.Context, payment *domain.Payment) error
}

// PaymentCreator handles payment creation business logic
type PaymentCreatorService struct {
	paymentStorer    PaymentStorer    // PaymentStorer implements the PaymentStorer interface
	walletReserver   WalletReserver   // WalletReserver implements the WalletReserver interface
	paymentPublisher PaymentPublisher // PaymentPublisher implements the PaymentPublisher interface
}

// NewPaymentCreator creates a new PaymentCreator
func NewPaymentCreatorService(ps PaymentStorer, wr WalletReserver, pp PaymentPublisher) (*PaymentCreatorService, error) {
	if ps == nil {
		return nil, errors.New("payment creator: storer cannot be nil")
	}
	if wr == nil {
		return nil, errors.New("payment creator: wallet reserver cannot be nil")
	}
	if pp == nil {
		return nil, errors.New("payment creator: publisher cannot be nil")
	}

	return &PaymentCreatorService{
		paymentStorer:    ps,
		walletReserver:   wr,
		paymentPublisher: pp,
	}, nil
}

// Create creates a new payment and publishes it
// It creates a new payment, reserves funds, updates the status to "reserved" and publishes the payment
// It returns a new payment and an error if the payment cannot be created
func (pcs *PaymentCreatorService) Create(ctx context.Context, idempotencyKey string, pr *PaymentRequest) (*domain.Payment, error) {
	// Step 1: Check if payment already exists
	existingPayment, err := pcs.paymentStorer.GetByIDempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("payment creator: get by idempotency key: %w", err)
	}
	if existingPayment != nil {
		return existingPayment, nil
	}

	// Step 2: Create payment with status "pending"
	payment := NewPayment(idempotencyKey, pr.UserID, pr.Amount, pr.Currency)

	if err := pcs.paymentStorer.Save(ctx, payment); err != nil {
		return nil, fmt.Errorf("payment creator: save payment: %w", err)
	}

	// Step 3: Reserve funds in wallet
	if err := pcs.walletReserver.Reserve(ctx, pr.UserID, pr.Amount, payment.ID); err != nil {
		if err := pcs.paymentStorer.UpdateStatus(ctx, payment.ID, domain.StatusFailed, ""); err != nil {
			return nil, fmt.Errorf("payment creator: update status to failed: %w", err)
		}
		return nil, fmt.Errorf("payment creator: reserve funds: %w", err)
	}

	// Step 4: Update status to "reserved"
	if err := pcs.paymentStorer.UpdateStatus(ctx, payment.ID, domain.StatusReserved, ""); err != nil {
		return nil, fmt.Errorf("payment creator: update status to reserved: %w", err)
	}

	payment.UpdateStatus(domain.StatusReserved)

	// Step 6: Publish payment event
	if err := pcs.paymentPublisher.Publish(ctx, payment); err != nil {
		return nil, fmt.Errorf("payment creator: publish payment: %w", err)
	}

	return payment, nil
}
