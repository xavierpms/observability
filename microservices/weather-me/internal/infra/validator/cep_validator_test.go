package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCEPFormatValid(t *testing.T) {
	validator := NewCEPValidator()

	assert.True(t, validator.ValidateCEPFormat("29902555"))
	assert.True(t, validator.ValidateCEPFormat("01001000"))
}

func TestValidateCEPFormatInvalid(t *testing.T) {
	validator := NewCEPValidator()

	testCases := []string{
		"2990255",
		"299025555",
		"2990255a",
		"29902-555",
		"",
		"abcd2555",
	}

	for _, cep := range testCases {
		assert.False(t, validator.ValidateCEPFormat(cep))
	}
}
