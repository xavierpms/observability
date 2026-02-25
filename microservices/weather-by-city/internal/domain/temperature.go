package domain

// Temperature represents the temperature in different scales
type Temperature struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

// CEPData represents the data returned by the ViaCEP API
type CEPData struct {
	CEP     string
	City    string
	Region  string
	RawData map[string]interface{}
}

// TemperatureRepository defines the contract for fetching temperature data
type TemperatureRepository interface {
	GetTemperatureByCityName(cityName string) (*Temperature, error)
}

// CEPRepository defines the contract for fetching CEP data
type CEPRepository interface {
	GetCEPData(cep string) (*CEPData, error)
}

// TemperatureUseCase defines the contract for the use case related to temperature retrieval
type TemperatureUseCase interface {
	GetTemperatureByCEP(cep string) (*Temperature, error)
}

// CEPValidator defines the contract for validating CEP
type CEPValidator interface {
	ValidateCEPFormat(cep string) bool
}
