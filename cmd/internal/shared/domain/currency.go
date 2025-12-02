package domain

import "errors"

// Currency represents the supported currencies for payments
type Currency string

const (
	CurrencyUSD Currency = "USD" // United States Dollar (USD)
	CurrencyEUR Currency = "EUR" // Euro (EUR)
	CurrencyGBP Currency = "ARS" // Argentine Peso (ARS)
)

// String converts the currency to a string
// It returns the string representation of the currency
func (c Currency) String() string {
	return string(c)
}

// Validate validates the currency
// It returns an error if the currency is invalid
func (c Currency) Validate() error {
	if c != CurrencyUSD && c != CurrencyEUR && c != CurrencyGBP {
		return errors.New("invalid currency")
	}
	return nil
}
