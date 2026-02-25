package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCEPFormatValid(t *testing.T) {
	// Arrange
	validator := NewCEPValidator()

	// Act & Assert
	assert.True(t, validator.ValidateCEPFormat("32450000"))
	assert.True(t, validator.ValidateCEPFormat("01021200"))
	assert.True(t, validator.ValidateCEPFormat("00000000"))
}

func TestValidateCEPFormatInvalid(t *testing.T) {
	// Arrange
	validator := NewCEPValidator()

	testCases := []string{
		"3245000",   // Very short (7 digits)
		"324500000", // Very long (9 digits)
		"3245000a",  // Contains letter
		"32450-000", // Contains hyphen
		"",          // Empty
		"   ",       // Only spaces
		"abcdefgh",  // All letters
	}

	for _, cep := range testCases {
		assert.False(t, validator.ValidateCEPFormat(cep))
	}
}
