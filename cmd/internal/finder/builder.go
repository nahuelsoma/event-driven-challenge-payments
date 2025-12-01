package finder

import "fmt"

func Build(db paymentReaderDB) (*Handler, error) {
	pr, err := NewPaymentReaderRepository(db)
	if err != nil {
		return nil, fmt.Errorf("finder builder: %w", err)
	}

	pf, err := NewPaymentFinderService(pr)
	if err != nil {
		return nil, fmt.Errorf("finder builder: %w", err)
	}

	h, err := NewHandler(pf)
	if err != nil {
		return nil, fmt.Errorf("finder builder: %w", err)
	}

	return h, nil
}
