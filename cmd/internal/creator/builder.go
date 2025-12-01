package creator

import (
	"github.com/nahuelsoma/event-driven-challenge-payments/infrastructure/messagebroker"
)

func Build(dc paymentStorerDB, c interface{}, mbc *messagebroker.Connection) (*Handler, error) {
	ps, err := NewPaymentStorerRepository(dc)
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
			Exchange:   "payments",
			RoutingKey: "payments",
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
