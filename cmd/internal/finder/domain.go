package finder

import "errors"

// PaymentFilter represents the filter criteria for finding a payment
type PaymentFilter struct {
	PaymentID string `json:"payment_id"` // Payment ID to search for
}

// Validate validates the payment filter
// It returns an error if the filter is invalid
func (f *PaymentFilter) Validate() error {
	if f.PaymentID == "" {
		return errors.New("payment ID is required")
	}
	return nil
}
