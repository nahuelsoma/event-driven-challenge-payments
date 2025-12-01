package processor

import (
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/http"
)

func Build(db paymentStorerDB, walletClient interface{}) (*Handler, error) {
	ps, err := NewPaymentStorerRepository(db)
	if err != nil {
		return nil, fmt.Errorf("processor builder: %w", err)
	}

	wc, err := NewWalletConfirmerRepository(walletClient)
	if err != nil {
		return nil, fmt.Errorf("processor builder: %w", err)
	}

	wr, err := NewWalletReleaserRepository(walletClient)
	if err != nil {
		return nil, fmt.Errorf("processor builder: %w", err)
	}

	// Create gateway client
	gatewayConfig := map[string]string{
		"host": "localhost",
		"port": "3000",
	}

	gatewayClient, err := http.NewHTTPClient(gatewayConfig)
	if err != nil {
		return nil, fmt.Errorf("processor builder: %w", err)
	}

	gp, err := NewGatewayProcessorRepository(gatewayClient)
	if err != nil {
		return nil, fmt.Errorf("processor builder: %w", err)
	}

	pps, err := NewPaymentProcessorService(ps, ps, wc, wr, gp)
	if err != nil {
		return nil, fmt.Errorf("processor builder: %w", err)
	}

	h, err := NewHandler(pps)
	if err != nil {
		return nil, fmt.Errorf("processor builder: %w", err)
	}

	return h, nil
}
