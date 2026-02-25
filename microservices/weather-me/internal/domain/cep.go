package domain

import "context"

// ServiceBResponse represents the HTTP response from weather-by-city
// that should be returned by weather-me.
type ServiceBResponse struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

// ServiceBRepository defines the contract for forwarding CEP requests to weather-by-city.
type ServiceBRepository interface {
	ForwardCEP(ctx context.Context, cep string) (*ServiceBResponse, error)
}

// CEPInputUseCase defines the contract for the CEP forwarding use case.
type CEPInputUseCase interface {
	ForwardCEP(ctx context.Context, cep string) (*ServiceBResponse, error)
}

// CEPValidator defines the contract for validating CEP format.
type CEPValidator interface {
	ValidateCEPFormat(cep string) bool
}
