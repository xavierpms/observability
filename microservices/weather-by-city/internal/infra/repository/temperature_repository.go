package repository

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/xavierpms/weather-by-city/internal/domain"
)

// WeatherAPIResponse represents the response from the WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

// TemperatureRepository implements domain.TemperatureRepository
type TemperatureRepository struct {
	apiURL string
	apiKey string
}

// NewTemperatureRepository creates a new temperature repository
func NewTemperatureRepository(apiURL, apiKey string) domain.TemperatureRepository {
	return &TemperatureRepository{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// GetTemperatureByCityName fetches the temperature for a given city
func (r *TemperatureRepository) GetTemperatureByCityName(cityName string) (*domain.Temperature, error) {
	// Encode the city name for URL
	encodedCityName := url.QueryEscape(cityName)

	// Build the URL with parameters
	requestURL := r.apiURL + "?q=" + encodedCityName + "&lang=pt&country=Brazil&key=" + r.apiKey
	log.Printf("calling Weather API: base_url=%s city=%s", requestURL, cityName)

	// Make the request
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Weather API request error: city=%s err=%v", cityName, err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("Weather API response: city=%s status=%d", cityName, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Weather API read response error: city=%s err=%v", cityName, err)
		return nil, err
	}

	// Unmarshal the response
	var weatherResp WeatherAPIResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		log.Printf("Weather API parse response error: city=%s err=%v", cityName, err)
		return nil, err
	}

	// Calculate the temperature in Kelvin
	kelvin := weatherResp.Current.TempC + 273.0
	log.Printf("Weather API request succeeded: city=%s temp_c=%.2f", cityName, weatherResp.Current.TempC)

	return &domain.Temperature{
		Celsius:    weatherResp.Current.TempC,
		Fahrenheit: weatherResp.Current.TempF,
		Kelvin:     kelvin,
	}, nil
}
