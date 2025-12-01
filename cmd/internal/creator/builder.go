package creator

import (
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/paymentstorer"
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

func Build(db paymentstorer.PaymentDB, c interface{}, mbc *messagebroker.Connection, exchange, queueName string) (*Handler, error) {
	ps, err := paymentstorer.NewStorer(db)
	if err != nil {
		return nil, err
	}

	wr, err := NewWalletReserverRepository(c)
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

	pc, err := NewPaymentCreatorService(ps, wr, pp)
	if err != nil {
		return nil, err
	}

	h, err := NewHandler(pc)
	if err != nil {
		return nil, err
	}

	return h, nil
}
