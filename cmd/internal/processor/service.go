package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/domain"
)

// PaymentResolver interface for resolving payment status
type PaymentResolver interface {
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
	UpdateStatus(ctx context.Context, paymentID string, status domain.Status, gatewayRef string) error
}

// WalletResolver interface for resolving funds
type WalletResolver interface {
	Confirm(ctx context.Context, userID string, amount float64, paymentID string) error
	Release(ctx context.Context, userID string, amount float64, paymentID string) error
}

// GatewayProcessor interface for processing payments with external gateway
type GatewayProcessor interface {
	Process(ctx context.Context, paymentID string, amount float64) (string, error)
}

// PaymentProcessorService handles payment processing business logic
type PaymentProcessorService struct {
	paymentResolver  PaymentResolver
	walletResolver   WalletResolver
	gatewayProcessor GatewayProcessor
}

// NewPaymentProcessorService creates a new PaymentProcessorService
func NewPaymentProcessorService(
	pr PaymentResolver,
	wr WalletResolver,
	gp GatewayProcessor,
) (*PaymentProcessorService, error) {
	if pr == nil {
		return nil, errors.New("payment processor: resolver cannot be nil")
	}
	if wr == nil {
		return nil, errors.New("payment processor: wallet resolver cannot be nil")
	}
	if gp == nil {
		return nil, errors.New("payment processor: gateway processor cannot be nil")
	}

	return &PaymentProcessorService{
		paymentResolver:  pr,
		walletResolver:   wr,
		gatewayProcessor: gp,
	}, nil
}

// Process processes a payment
// Flow: Check status (idempotency) → Process with gateway → Confirm/Release funds → Update status
func (pps *PaymentProcessorService) Process(ctx context.Context, payment *domain.Payment) error {
	// Step 1: Check payment status for idempotency
	existing, err := pps.paymentResolver.GetByID(ctx, payment.ID)
	if err != nil {
		return fmt.Errorf("payment processor: failed to get payment: %w", err)
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
		if releaseErr := pps.walletResolver.Release(ctx, payment.UserID, payment.Amount, payment.ID); releaseErr != nil {
			return fmt.Errorf("payment processor: failed to release funds: %w", releaseErr)
		}

		if updateErr := pps.paymentResolver.UpdateStatus(ctx, payment.ID, domain.StatusFailed, ""); updateErr != nil {
			return fmt.Errorf("payment processor: failed to update status to failed: %w", updateErr)
		}

		return nil // Payment failed but handled correctly
	}

	// Step 3: Gateway succeeded → Confirm funds
	if err := pps.walletResolver.Confirm(ctx, payment.UserID, payment.Amount, payment.ID); err != nil {
		return fmt.Errorf("payment processor: failed to confirm funds: %w", err)
	}

	// Step 4: Update status to completed
	if err := pps.paymentResolver.UpdateStatus(ctx, payment.ID, domain.StatusCompleted, gatewayRef); err != nil {
		return fmt.Errorf("payment processor: failed to update status to completed: %w", err)
	}

	return nil
}
