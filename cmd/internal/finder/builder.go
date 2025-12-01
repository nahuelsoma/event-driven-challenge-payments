package finder

func Build(db paymentReaderDB) (*Handler, error) {
	pr, err := NewPaymentReaderRepository(db)
	if err != nil {
		return nil, err
	}

	pf, err := NewPaymentFinderService(pr)
	if err != nil {
		return nil, err
	}

	h, err := NewHandler(pf)
	if err != nil {
		return nil, err
	}

	return h, nil
}
