package validator

import "strconv"

// CEPValidatorImpl implements domain.CEPValidator.
type CEPValidatorImpl struct{}

// NewCEPValidator creates a new CEP validator.
func NewCEPValidator() *CEPValidatorImpl {
	return &CEPValidatorImpl{}
}

// ValidateCEPFormat validates if CEP has exactly 8 numeric digits.
func (v *CEPValidatorImpl) ValidateCEPFormat(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	_, err := strconv.Atoi(cep)
	return err == nil
}
