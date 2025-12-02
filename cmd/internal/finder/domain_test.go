package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentFilter_Validate(t *testing.T) {
	tests := []struct {
		name          string
		filter        *PaymentFilter
		expectedError string
	}{
		{
			name: "when payment ID is provided it should pass validation and no error",
			filter: &PaymentFilter{
				PaymentID: "pay_123",
			},
			expectedError: "",
		},
		{
			name: "when payment ID is empty it should return error with message 'payment ID is required'",
			filter: &PaymentFilter{
				PaymentID: "",
			},
			expectedError: "payment ID is required",
		},
		{
			name: "when payment ID has minimum valid length it should pass validation and no error",
			filter: &PaymentFilter{
				PaymentID: "p",
			},
			expectedError: "",
		},
		{
			name: "when payment ID is a long string it should pass validation and no error",
			filter: &PaymentFilter{
				PaymentID: "pay_very_long_id_with_many_characters_123456789",
			},
			expectedError: "",
		},
		{
			name: "when payment ID contains special characters it should pass validation and no error",
			filter: &PaymentFilter{
				PaymentID: "pay_123-abc_456",
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// (Filter already prepared in test struct)

			// Act
			err := tt.filter.Validate()

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
