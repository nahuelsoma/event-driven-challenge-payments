package finder

import (
	"github.com/nahuelsoma/event-driven-challenge-payments/cmd/internal/shared/repository/paymentstorer"
)

func Build(db paymentstorer.PaymentDB) (*Handler, error) {
	ps, err := paymentstorer.NewStorer(db)
	if err != nil {
		return nil, err
	}

	pf, err := NewPaymentFinderService(ps)
	if err != nil {
		return nil, err
	}

	h, err := NewHandler(pf)
	if err != nil {
		return nil, err
	}

	return h, nil
}
