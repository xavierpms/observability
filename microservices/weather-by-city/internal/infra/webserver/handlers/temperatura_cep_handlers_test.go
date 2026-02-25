package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/xavierpms/weather-by-city/internal/domain"
)

// MockTemperatureUseCase is a mock of the TemperatureUseCase for testing
type MockTemperatureUseCase struct {
	getTemperatureByCEPFunc func(ctx context.Context, cep string) (*domain.Temperature, error)
	receivedCEP             string
}

func (m *MockTemperatureUseCase) GetTemperatureByCEP(ctx context.Context, cep string) (*domain.Temperature, error) {
	m.receivedCEP = cep
	return m.getTemperatureByCEPFunc(ctx, cep)
}

func buildRequestWithCEP(cep string) *http.Request {
	req := httptest.NewRequest("GET", "/"+cep, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", cep)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

// TestGetTemperatureByCEPSuccess tests the success of the request
func TestGetTemperatureByCEPSuccess(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(ctx context.Context, cep string) (*domain.Temperature, error) {
			return &domain.Temperature{
				City:       "São Paulo",
				Celsius:    28.5,
				Fahrenheit: 83.3,
				Kelvin:     301.65,
			}, nil
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := buildRequestWithCEP("32450000")
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var temperature domain.Temperature
	err := json.Unmarshal(w.Body.Bytes(), &temperature)
	assert.NoError(t, err)
	assert.Equal(t, "32450000", mockUseCase.receivedCEP)
	assert.Equal(t, "São Paulo", temperature.City)
	assert.Equal(t, 28.5, temperature.Celsius)
	assert.Equal(t, 83.3, temperature.Fahrenheit)
	assert.Equal(t, 301.65, temperature.Kelvin)
}

// TestGetTemperatureByCEPInvalidFormat tests the case when the CEP has an invalid format
func TestGetTemperatureByCEPInvalidFormat(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(ctx context.Context, cep string) (*domain.Temperature, error) {
			return nil, domain.ErrInvalidCEPFormat
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := buildRequestWithCEP("3245000000")
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var errResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid zipcode", errResponse.Message)
}

// TestGetTemperatureByCEPNotFound tests the case when the CEP is not found
func TestGetTemperatureByCEPNotFound(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(ctx context.Context, cep string) (*domain.Temperature, error) {
			return nil, domain.ErrCEPNotFound
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := buildRequestWithCEP("00000000")
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Cannot find zipcode", errResponse.Message)
}

// TestGetTemperatureByCEPTemperatureNotFound tests the case when the temperature is not found
func TestGetTemperatureByCEPTemperatureNotFound(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(ctx context.Context, cep string) (*domain.Temperature, error) {
			return nil, domain.ErrTemperatureNotFound
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := buildRequestWithCEP("32450000")
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Cannot fetch temperature", errResponse.Message)
}

func TestGetTemperatureByCEPInternalError(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(ctx context.Context, cep string) (*domain.Temperature, error) {
			return nil, errors.New("unexpected")
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := buildRequestWithCEP("32450000")
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Internal server error", errResponse.Message)
}
