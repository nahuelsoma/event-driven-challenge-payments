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
// It returns a new GatewayProcessorRepository and an error if the client is nil
func NewGatewayProcessorRepository(client *http.Client) (*GatewayProcessorRepository, error) {
	if client == nil {
		return nil, errors.New("gateway processor: client cannot be nil")
	}

	return &GatewayProcessorRepository{client: client}, nil
}

// Process processes a payment with the external gateway
// It processes a payment with the external gateway and returns the gateway reference on success
func (r *GatewayProcessorRepository) Process(ctx context.Context, paymentID string, amount float64) (string, error) {
	slog.DebugContext(ctx, "[DEBUG] GatewayProcessorRepository.Process called", "payment_id", paymentID, "amount", amount)
	// TODO: Implement the logic to process the payment via external gateway
	// For now, return a mock gateway reference
	gatewayRef := "gw_" + uuid.New().String()
	return gatewayRef, nil
}
