package processor

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// GatewayProcessorRepository processes payments with external gateway
type GatewayProcessorRepository struct {
	client *http.Client
}

// NewGatewayProcessorRepository creates a new GatewayProcessorRepository
func NewGatewayProcessorRepository(client *http.Client) (*GatewayProcessorRepository, error) {
	if client == nil {
		return nil, errors.New("gateway processor: client cannot be nil")
	}

	return &GatewayProcessorRepository{client: client}, nil
}

// Process processes a payment with the external gateway
// Returns the gateway reference on success
func (r *GatewayProcessorRepository) Process(ctx context.Context, paymentID string, amount float64) (string, error) {
	slog.DebugContext(ctx, "[DEBUG] GatewayProcessorRepository.Process called", "payment_id", paymentID, "amount", amount)
	// TODO: Implement the logic to process the payment via external gateway
	// For now, return a mock gateway reference
	gatewayRef := "gw_" + uuid.New().String()
	return gatewayRef, nil
}
