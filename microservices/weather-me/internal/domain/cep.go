package domain

// ServiceBResponse represents the HTTP response from service-b
// that should be returned by service-a.
type ServiceBResponse struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

// ServiceBRepository defines the contract for forwarding CEP requests to service-b.
type ServiceBRepository interface {
	ForwardCEP(cep string) (*ServiceBResponse, error)
}

// CEPInputUseCase defines the contract for the CEP forwarding use case.
type CEPInputUseCase interface {
	ForwardCEP(cep string) (*ServiceBResponse, error)
}

// CEPValidator defines the contract for validating CEP format.
type CEPValidator interface {
	ValidateCEPFormat(cep string) bool
}
