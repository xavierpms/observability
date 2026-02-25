package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xavierpms/weather-by-city/internal/domain"
)

type mockCEPRepository struct {
	getCEPDataFunc func(cep string) (*domain.CEPData, error)
	called         bool
	receivedCEP    string
}

func (m *mockCEPRepository) GetCEPData(cep string) (*domain.CEPData, error) {
	m.called = true
	m.receivedCEP = cep
	return m.getCEPDataFunc(cep)
}

type mockTemperatureRepository struct {
	getTemperatureByCityNameFunc func(cityName string) (*domain.Temperature, error)
	called                       bool
	receivedCity                 string
}

func (m *mockTemperatureRepository) GetTemperatureByCityName(cityName string) (*domain.Temperature, error) {
	m.called = true
	m.receivedCity = cityName
	return m.getTemperatureByCityNameFunc(cityName)
}

type mockCEPValidator struct {
	validateCEPFormatFunc func(cep string) bool
}

func (m *mockCEPValidator) ValidateCEPFormat(cep string) bool {
	return m.validateCEPFormatFunc(cep)
}

func TestGetTemperatureByCEPInvalidFormat(t *testing.T) {
	cepRepo := &mockCEPRepository{getCEPDataFunc: func(cep string) (*domain.CEPData, error) {
		return &domain.CEPData{}, nil
	}}
	tempRepo := &mockTemperatureRepository{getTemperatureByCityNameFunc: func(cityName string) (*domain.Temperature, error) {
		return &domain.Temperature{}, nil
	}}
	validator := &mockCEPValidator{validateCEPFormatFunc: func(cep string) bool {
		return false
	}}

	useCase := NewGetTemperatureByCEP(cepRepo, tempRepo, validator)
	temperature, err := useCase.GetTemperatureByCEP("3245000")

	assert.Nil(t, temperature)
	assert.ErrorIs(t, err, domain.ErrInvalidCEPFormat)
	assert.False(t, cepRepo.called)
	assert.False(t, tempRepo.called)
}

func TestGetTemperatureByCEPCEPNotFound(t *testing.T) {
	cepRepo := &mockCEPRepository{getCEPDataFunc: func(cep string) (*domain.CEPData, error) {
		return nil, errors.New("viacep error")
	}}
	tempRepo := &mockTemperatureRepository{getTemperatureByCityNameFunc: func(cityName string) (*domain.Temperature, error) {
		return &domain.Temperature{}, nil
	}}
	validator := &mockCEPValidator{validateCEPFormatFunc: func(cep string) bool {
		return true
	}}

	useCase := NewGetTemperatureByCEP(cepRepo, tempRepo, validator)
	temperature, err := useCase.GetTemperatureByCEP("32450000")

	assert.Nil(t, temperature)
	assert.ErrorIs(t, err, domain.ErrCEPNotFound)
	assert.True(t, cepRepo.called)
	assert.Equal(t, "32450000", cepRepo.receivedCEP)
	assert.False(t, tempRepo.called)
}

func TestGetTemperatureByCEPTemperatureNotFound(t *testing.T) {
	cepRepo := &mockCEPRepository{getCEPDataFunc: func(cep string) (*domain.CEPData, error) {
		return &domain.CEPData{City: "São Paulo"}, nil
	}}
	tempRepo := &mockTemperatureRepository{getTemperatureByCityNameFunc: func(cityName string) (*domain.Temperature, error) {
		return nil, errors.New("weather error")
	}}
	validator := &mockCEPValidator{validateCEPFormatFunc: func(cep string) bool {
		return true
	}}

	useCase := NewGetTemperatureByCEP(cepRepo, tempRepo, validator)
	temperature, err := useCase.GetTemperatureByCEP("32450000")

	assert.Nil(t, temperature)
	assert.ErrorIs(t, err, domain.ErrTemperatureNotFound)
	assert.True(t, tempRepo.called)
	assert.Equal(t, "São Paulo", tempRepo.receivedCity)
}

func TestGetTemperatureByCEPSuccess(t *testing.T) {
	cepRepo := &mockCEPRepository{getCEPDataFunc: func(cep string) (*domain.CEPData, error) {
		return &domain.CEPData{City: "São Paulo"}, nil
	}}
	tempRepo := &mockTemperatureRepository{getTemperatureByCityNameFunc: func(cityName string) (*domain.Temperature, error) {
		return &domain.Temperature{
			Celsius:    28.5,
			Fahrenheit: 83.3,
			Kelvin:     301.5,
		}, nil
	}}
	validator := &mockCEPValidator{validateCEPFormatFunc: func(cep string) bool {
		return true
	}}

	useCase := NewGetTemperatureByCEP(cepRepo, tempRepo, validator)
	temperature, err := useCase.GetTemperatureByCEP("32450000")

	assert.NoError(t, err)
	assert.NotNil(t, temperature)
	assert.Equal(t, "São Paulo", temperature.City)
	assert.Equal(t, 28.5, temperature.Celsius)
	assert.Equal(t, 83.3, temperature.Fahrenheit)
	assert.Equal(t, 301.5, temperature.Kelvin)
	assert.True(t, cepRepo.called)
	assert.True(t, tempRepo.called)
	assert.Equal(t, "São Paulo", tempRepo.receivedCity)
}
