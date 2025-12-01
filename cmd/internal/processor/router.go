package processor

// Start initializes and returns the processor handler
// The handler implements the MessageHandler interface from infrastructure/messagebroker
func Start(db paymentStorerDB, walletClient interface{}, gatewayClient interface{}) (*Handler, error) {
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

	gp, err := NewGatewayProcessorRepository(gatewayClient)
	if err != nil {
		return nil, err
	}

	// PaymentStorerRepository implements both PaymentReader and PaymentUpdater interfaces
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
