package repository

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/xavierpms/weather-by-city/internal/domain"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	client *http.Client
}

// NewTemperatureRepository creates a new temperature repository
func NewTemperatureRepository(apiURL, apiKey string) domain.TemperatureRepository {
	return &TemperatureRepository{
		apiURL: apiURL,
		apiKey: apiKey,
		client: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

// GetTemperatureByCityName fetches the temperature for a given city
func (r *TemperatureRepository) GetTemperatureByCityName(ctx context.Context, cityName string) (*domain.Temperature, error) {
	tracer := otel.Tracer("weather-by-city.repository.weatherapi")
	ctx, span := tracer.Start(ctx, "weatherapi.lookup")
	defer span.End()

	// Encode the city name for URL
	encodedCityName := url.QueryEscape(cityName)

	// Build the URL with parameters
	requestURL := r.apiURL + "?q=" + encodedCityName + "&lang=pt&country=Brazil&key=" + r.apiKey
	span.SetAttributes(attribute.String("http.url", requestURL), attribute.String("city", cityName))
	log.Printf("calling Weather API: base_url=%s city=%s", requestURL, cityName)

	// Make the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to build Weather API request")
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		log.Printf("Weather API request error: city=%s err=%v", cityName, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Weather API request failed")
		return nil, err
	}
	defer resp.Body.Close()
	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	log.Printf("Weather API response: city=%s status=%d", cityName, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Weather API read response error: city=%s err=%v", cityName, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed reading Weather API response")
		return nil, err
	}

	// Unmarshal the response
	var weatherResp WeatherAPIResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		log.Printf("Weather API parse response error: city=%s err=%v", cityName, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed parsing Weather API response")
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
