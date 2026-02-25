package repository

import (
	"io"
	"net/http"
	"strings"

	"github.com/xavierpms/service-a/internal/domain"
)

// ServiceBRepositoryImpl implements domain.ServiceBRepository.
type ServiceBRepositoryImpl struct {
	apiURL string
}

// NewServiceBRepository creates a new repository for service-b HTTP calls.
func NewServiceBRepository(apiURL string) domain.ServiceBRepository {
	return &ServiceBRepositoryImpl{apiURL: strings.TrimRight(apiURL, "/")}
}

// ForwardCEP forwards CEP to service-b and returns raw response details.
func (r *ServiceBRepositoryImpl) ForwardCEP(cep string) (*domain.ServiceBResponse, error) {
	requestURL := r.apiURL + "/" + cep

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &domain.ServiceBResponse{
		StatusCode:  resp.StatusCode,
		Body:        body,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
}
