package processor

import (
	"log"

	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/http"
)

func Build(db paymentStorerDB, walletClient interface{}) (*Handler, error) {
	ps, err := NewPaymentStorerRepository(db)
	if err != nil {
		return nil, err
	}

	wc, err := NewWalletConfirmerRepository(walletClient)
	if err != nil {
		return nil, err
	}

	wr, err := NewWalletReleaserRepository(walletClient)
	if err != nil {
		return nil, err
	}

	// Create gateway client
	gatewayConfig := map[string]string{
		"host": "localhost",
		"port": "3000",
	}

	gatewayClient, err := http.NewHTTPClient(gatewayConfig)
	if err != nil {
		log.Fatalf("api: failed to create HTTP client: %v", err)
	}

	gp, err := NewGatewayProcessorRepository(gatewayClient)
	if err != nil {
		return nil, err
	}

	pps, err := NewPaymentProcessorService(ps, ps, wc, wr, gp)
	if err != nil {
		return nil, err
	}

	h, err := NewHandler(pps)
	if err != nil {
		return nil, err
	}

	return h, nil
}
