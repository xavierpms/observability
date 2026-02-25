package validator

import "strconv"

// CEPValidatorImpl implements domain.CEPValidator
type CEPValidatorImpl struct{}

// NewCEPValidator creates a new CEP validator
func NewCEPValidator() *CEPValidatorImpl {
	return &CEPValidatorImpl{}
}

// ValidateCEPFormat validates the format of the CEP (must have 8 digits)
func (v *CEPValidatorImpl) ValidateCEPFormat(cep string) bool {
	// Check if it has exactly 8 characters
	if len(cep) != 8 {
		return false
	}

	// Check if all are numbers
	_, err := strconv.Atoi(cep)
	return err == nil
}
