package usecase

import "github.com/xavierpms/service-a/internal/domain"

// ForwardCEPUseCase represents the use case for forwarding CEP to service-b.
type ForwardCEPUseCase struct {
	serviceBRepository domain.ServiceBRepository
	cepValidator       domain.CEPValidator
}

// NewForwardCEPUseCase creates a new use case instance.
func NewForwardCEPUseCase(
	serviceBRepo domain.ServiceBRepository,
	validator domain.CEPValidator,
) domain.CEPInputUseCase {
	return &ForwardCEPUseCase{
		serviceBRepository: serviceBRepo,
		cepValidator:       validator,
	}
}

// ForwardCEP validates CEP and forwards it to service-b.
func (u *ForwardCEPUseCase) ForwardCEP(cep string) (*domain.ServiceBResponse, error) {
	if !u.cepValidator.ValidateCEPFormat(cep) {
		return nil, domain.ErrInvalidCEPFormat
	}

	response, err := u.serviceBRepository.ForwardCEP(cep)
	if err != nil {
		return nil, domain.ErrForwardCEP
	}

	return response, nil
}
