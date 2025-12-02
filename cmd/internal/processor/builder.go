package processor

import (
	"log"
	"time"

	"net/http"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/walletclient"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/restclient"
)

// Build creates a new Handler with all dependencies wired up
// It builds a new handler and returns an error if the builder fails
func Build(db paymentstorer.PaymentDB, rc *http.Client) (*Handler, error) {
	ps, err := paymentstorer.NewStorer(db)
	if err != nil {
		return nil, err
	}

	wc, err := walletclient.NewWalletClient(rc)
	if err != nil {
		return nil, err
	}

	gatewayConfig := &restclient.Config{
		BaseURL: "http://test-gateway-service.com",
		Timeout: 300 * time.Millisecond,
	}

	gatewayClient, err := restclient.NewRestClient(gatewayConfig)
	if err != nil {
		log.Fatalf("api: failed to create HTTP client: %v", err)
	}

	gpr, err := NewGatewayProcessorRepository(gatewayClient)
	if err != nil {
		return nil, err
	}

	pps, err := NewPaymentProcessorService(ps, wc, gpr)
	if err != nil {
		return nil, err
	}

	h, err := NewHandler(pps)
	if err != nil {
		return nil, err
	}

	return h, nil
}
