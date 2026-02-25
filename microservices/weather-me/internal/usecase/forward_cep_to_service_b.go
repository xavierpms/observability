package usecase

import (
	"context"

	"github.com/xavierpms/service-a/internal/domain"
)

// ForwardCEPUseCase represents the use case for forwarding CEP to weather-by-city.
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

// ForwardCEP validates CEP and forwards it to weather-by-city.
func (u *ForwardCEPUseCase) ForwardCEP(ctx context.Context, cep string) (*domain.ServiceBResponse, error) {
	if !u.cepValidator.ValidateCEPFormat(cep) {
		return nil, domain.ErrInvalidCEPFormat
	}

	response, err := u.serviceBRepository.ForwardCEP(ctx, cep)
	if err != nil {
		return nil, domain.ErrForwardCEP
	}

	return response, nil
}
