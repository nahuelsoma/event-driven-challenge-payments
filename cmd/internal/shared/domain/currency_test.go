package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency_String(t *testing.T) {
	tests := []struct {
		name           string
		currency       Currency
		expectedResult string
	}{
		{
			name:           "when currency is USD it should return USD string and no error",
			currency:       CurrencyUSD,
			expectedResult: "USD",
		},
		{
			name:           "when currency is EUR it should return EUR string and no error",
			currency:       CurrencyEUR,
			expectedResult: "EUR",
		},
		{
			name:           "when currency is GBP it should return ARS string and no error",
			currency:       CurrencyGBP,
			expectedResult: "ARS",
		},
		{
			name:           "when currency is custom value it should return that value as string and no error",
			currency:       Currency("BRL"),
			expectedResult: "BRL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Currency already prepared in test struct)

			// Act
			result := tt.currency.String()

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestCurrency_Validate(t *testing.T) {
	tests := []struct {
		name          string
		currency      Currency
		expectedError error
	}{
		{
			name:          "when currency is USD it should pass validation and no error",
			currency:      CurrencyUSD,
			expectedError: nil,
		},
		{
			name:          "when currency is EUR it should pass validation and no error",
			currency:      CurrencyEUR,
			expectedError: nil,
		},
		{
			name:          "when currency is GBP it should pass validation and no error",
			currency:      CurrencyGBP,
			expectedError: nil,
		},
		{
			name:          "when currency is invalid it should return error with message containing 'invalid currency'",
			currency:      Currency("INVALID"),
			expectedError: assert.AnError,
		},
		{
			name:          "when currency is empty it should return error with message containing 'invalid currency'",
			currency:      Currency(""),
			expectedError: assert.AnError,
		},
		{
			name:          "when currency is lowercase usd it should return error with message containing 'invalid currency'",
			currency:      Currency("usd"),
			expectedError: assert.AnError,
		},
		{
			name:          "when currency is BRL it should return error with message containing 'invalid currency'",
			currency:      Currency("BRL"),
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Currency already prepared in test struct)

			// Act
			err := tt.currency.Validate()

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, "invalid currency", err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
