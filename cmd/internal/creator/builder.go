package creator

import (
	"fmt"

	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

func Build(dc paymentStorerDB, c interface{}, mbc *messagebroker.Connection, exchange, queueName string) (*Handler, error) {
	ps, err := NewPaymentStorerRepository(dc)
	if err != nil {
		return nil, fmt.Errorf("creator builder: %w", err)
	}

	wr, err := NewWalletReserverRepository(c)
	if err != nil {
		return nil, fmt.Errorf("creator builder: %w", err)
	}

	p, err := messagebroker.NewPublisher(
		mbc,
		messagebroker.PublisherConfig{
			Exchange:   exchange,
			RoutingKey: queueName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creator builder: %w", err)
	}

	pp, err := NewPaymentPublisherRepository(p)
	if err != nil {
		return nil, fmt.Errorf("creator builder: %w", err)
	}

	pc, err := NewPaymentCreatorService(ps, wr, pp)
	if err != nil {
		return nil, fmt.Errorf("creator builder: %w", err)
	}

	h, err := NewHandler(pc)
	if err != nil {
		return nil, fmt.Errorf("creator builder: %w", err)
	}

	return h, nil
}
