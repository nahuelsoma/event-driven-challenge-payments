package creator

import (
	"net/http"

	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/walletclient"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

// Build creates a new Handler with all dependencies wired up
func Build(db paymentstorer.PaymentDB, rc *http.Client, mbc *messagebroker.Connection, exchange, queueName string) (*Handler, error) {
	ps, err := paymentstorer.NewStorer(db)
	if err != nil {
		return nil, err
	}

	wc, err := walletclient.NewWalletClient(rc)
	if err != nil {
		return nil, err
	}

	p, err := messagebroker.NewPublisher(
		mbc,
		messagebroker.PublisherConfig{
			Exchange:   exchange,  // Topic exchange for routing
			RoutingKey: queueName, // Queue name as routing key (e.g., payments.created)
		},
	)
	if err != nil {
		return nil, err
	}

	pp, err := NewPaymentPublisherRepository(p)
	if err != nil {
		return nil, err
	}

	pc, err := NewPaymentCreatorService(ps, wc, pp)
	if err != nil {
		return nil, err
	}

	h, err := NewHandler(pc)
	if err != nil {
		return nil, err
	}

	return h, nil
}
