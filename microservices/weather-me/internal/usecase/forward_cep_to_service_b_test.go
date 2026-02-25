package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xavierpms/service-a/internal/domain"
)

type mockServiceBRepository struct {
	forwardCEPFunc func(ctx context.Context, cep string) (*domain.ServiceBResponse, error)
	called         bool
	receivedCEP    string
}

func (m *mockServiceBRepository) ForwardCEP(ctx context.Context, cep string) (*domain.ServiceBResponse, error) {
	m.called = true
	m.receivedCEP = cep
	return m.forwardCEPFunc(ctx, cep)
}

type mockCEPValidator struct {
	validateCEPFormatFunc func(cep string) bool
}

func (m *mockCEPValidator) ValidateCEPFormat(cep string) bool {
	return m.validateCEPFormatFunc(cep)
}

func TestForwardCEPInvalidFormat(t *testing.T) {
	repo := &mockServiceBRepository{forwardCEPFunc: func(ctx context.Context, cep string) (*domain.ServiceBResponse, error) {
		return &domain.ServiceBResponse{}, nil
	}}
	validator := &mockCEPValidator{validateCEPFormatFunc: func(cep string) bool {
		return false
	}}

	useCase := NewForwardCEPUseCase(repo, validator)
	response, err := useCase.ForwardCEP(context.Background(), "2990255")

	assert.Nil(t, response)
	assert.ErrorIs(t, err, domain.ErrInvalidCEPFormat)
	assert.False(t, repo.called)
}

func TestForwardCEPRepositoryError(t *testing.T) {
	repo := &mockServiceBRepository{forwardCEPFunc: func(ctx context.Context, cep string) (*domain.ServiceBResponse, error) {
		return nil, errors.New("weather-by-city service is down")
	}}
	validator := &mockCEPValidator{validateCEPFormatFunc: func(cep string) bool {
		return true
	}}

	useCase := NewForwardCEPUseCase(repo, validator)
	response, err := useCase.ForwardCEP(context.Background(), "29902555")

	assert.Nil(t, response)
	assert.ErrorIs(t, err, domain.ErrForwardCEP)
	assert.True(t, repo.called)
	assert.Equal(t, "29902555", repo.receivedCEP)
}

func TestForwardCEPSuccess(t *testing.T) {
	repo := &mockServiceBRepository{forwardCEPFunc: func(ctx context.Context, cep string) (*domain.ServiceBResponse, error) {
		return &domain.ServiceBResponse{
			StatusCode:  200,
			Body:        []byte(`{"city":"Vitória"}`),
			ContentType: "application/json",
		}, nil
	}}
	validator := &mockCEPValidator{validateCEPFormatFunc: func(cep string) bool {
		return true
	}}

	useCase := NewForwardCEPUseCase(repo, validator)
	response, err := useCase.ForwardCEP(context.Background(), "29902555")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "application/json", response.ContentType)
	assert.Equal(t, []byte(`{"city":"Vitória"}`), response.Body)
	assert.Equal(t, "29902555", repo.receivedCEP)
}
