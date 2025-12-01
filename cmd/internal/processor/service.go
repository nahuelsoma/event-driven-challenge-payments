package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// PaymentReader interface for reading payment status
type PaymentReader interface {
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
}

// PaymentUpdater interface for updating payment status
type PaymentUpdater interface {
	UpdateStatus(ctx context.Context, paymentID string, status domain.Status, gatewayRef string) error
}

// WalletConfirmer interface for confirming funds
type WalletConfirmer interface {
	Confirm(ctx context.Context, userID string, amount float64, paymentID string) error
}

// WalletReleaser interface for releasing funds
type WalletReleaser interface {
	Release(ctx context.Context, userID string, amount float64, paymentID string) error
}

// GatewayProcessor interface for processing payments with external gateway
type GatewayProcessor interface {
	Process(ctx context.Context, paymentID string, amount float64) (string, error)
}

// PaymentProcessorService handles payment processing business logic
type PaymentProcessorService struct {
	paymentReader    PaymentReader
	paymentUpdater   PaymentUpdater
	walletConfirmer  WalletConfirmer
	walletReleaser   WalletReleaser
	gatewayProcessor GatewayProcessor
}

// NewPaymentProcessorService creates a new PaymentProcessorService
func NewPaymentProcessorService(
	pr PaymentReader,
	pu PaymentUpdater,
	wc WalletConfirmer,
	wr WalletReleaser,
	gp GatewayProcessor,
) (*PaymentProcessorService, error) {
	if pr == nil {
		return nil, errors.New("payment processor: reader cannot be nil")
	}
	if pu == nil {
		return nil, errors.New("payment processor: updater cannot be nil")
	}
	if wc == nil {
		return nil, errors.New("payment processor: wallet confirmer cannot be nil")
	}
	if wr == nil {
		return nil, errors.New("payment processor: wallet releaser cannot be nil")
	}
	if gp == nil {
		return nil, errors.New("payment processor: gateway processor cannot be nil")
	}

	return &PaymentProcessorService{
		paymentReader:    pr,
		paymentUpdater:   pu,
		walletConfirmer:  wc,
		walletReleaser:   wr,
		gatewayProcessor: gp,
	}, nil
}

// Process processes a payment
// Flow: Check status (idempotency) → Process with gateway → Confirm/Release funds → Update status
func (pps *PaymentProcessorService) Process(ctx context.Context, payment *domain.Payment) error {
	// Step 1: Check payment status for idempotency
	existing, err := pps.paymentReader.GetByID(ctx, payment.ID)
	if err != nil {
		return fmt.Errorf("payment processor: get payment: %w", err)
	}

	// Skip if already processed (idempotency)
	switch existing.Status {
	case domain.StatusCompleted, domain.StatusFailed:
		return nil // Already processed, skip silently
	case domain.StatusReserved:
		// Continue processing
	default:
		return nil // Unexpected status, skip
	}

	// Step 2: Process with gateway
	gatewayRef, err := pps.gatewayProcessor.Process(ctx, payment.ID, payment.Amount)
	if err != nil {
		// Gateway failed → Release funds and mark as failed
		if releaseErr := pps.walletReleaser.Release(ctx, payment.UserID, payment.Amount, payment.ID); releaseErr != nil {
			return fmt.Errorf("payment processor: release funds: %w", releaseErr)
		}

		if updateErr := pps.paymentUpdater.UpdateStatus(ctx, payment.ID, domain.StatusFailed, ""); updateErr != nil {
			return fmt.Errorf("payment processor: update status failed: %w", updateErr)
		}

		return nil // Payment failed but handled correctly
	}

	// Step 3: Gateway succeeded → Confirm funds
	if err := pps.walletConfirmer.Confirm(ctx, payment.UserID, payment.Amount, payment.ID); err != nil {
		return fmt.Errorf("payment processor: confirm funds: %w", err)
	}

	// Step 4: Update status to completed
	if err := pps.paymentUpdater.UpdateStatus(ctx, payment.ID, domain.StatusCompleted, gatewayRef); err != nil {
		return fmt.Errorf("payment processor: update status completed: %w", err)
	}

	return nil
}
