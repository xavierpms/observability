package usecase

import (
	"github.com/xavierpms/weather-by-city/internal/domain"
)

// GetTemperatureByCEP represents the use case for fetching temperature by CEP
type GetTemperatureByCEP struct {
	cepRepository         domain.CEPRepository
	temperatureRepository domain.TemperatureRepository
	cepValidator          domain.CEPValidator
}

// NewGetTemperatureByCEP creates a new instance of the use case
func NewGetTemperatureByCEP(
	cepRepo domain.CEPRepository,
	tempRepo domain.TemperatureRepository,
	validator domain.CEPValidator,
) domain.TemperatureUseCase {
	return &GetTemperatureByCEP{
		cepRepository:         cepRepo,
		temperatureRepository: tempRepo,
		cepValidator:          validator,
	}
}

// GetTemperatureByCEP executes the business logic
func (u *GetTemperatureByCEP) GetTemperatureByCEP(cep string) (*domain.Temperature, error) {
	// Validate the CEP format
	if !u.cepValidator.ValidateCEPFormat(cep) {
		return nil, domain.ErrInvalidCEPFormat
	}

	// Fetch the CEP data
	cepData, err := u.cepRepository.GetCEPData(cep)
	if err != nil {
		return nil, domain.ErrCEPNotFound
	}

	// Fetch the temperature for the city
	temperature, err := u.temperatureRepository.GetTemperatureByCityName(cepData.City)
	if err != nil {
		return nil, domain.ErrTemperatureNotFound
	}

	temperature.City = cepData.City

	return temperature, nil
}
