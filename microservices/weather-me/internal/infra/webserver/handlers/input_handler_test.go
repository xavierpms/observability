package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xavierpms/service-a/internal/domain"
)

type mockCEPInputUseCase struct {
	forwardCEPFunc func(cep string) (*domain.ServiceBResponse, error)
	receivedCEP    string
}

func (m *mockCEPInputUseCase) ForwardCEP(cep string) (*domain.ServiceBResponse, error) {
	m.receivedCEP = cep
	return m.forwardCEPFunc(cep)
}

func TestForwardCEPSuccess(t *testing.T) {
	mockUseCase := &mockCEPInputUseCase{
		forwardCEPFunc: func(cep string) (*domain.ServiceBResponse, error) {
			return &domain.ServiceBResponse{
				StatusCode:  200,
				Body:        []byte(`{"city":"Vitória","temp_C":28.0}`),
				ContentType: "application/json",
			}, nil
		},
	}

	handler := NewInputHandler(mockUseCase)
	body := []byte(`{"cep":"29902555"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ForwardCEP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "29902555", mockUseCase.receivedCEP)
	assert.JSONEq(t, `{"city":"Vitória","temp_C":28.0}`, w.Body.String())
}

func TestForwardCEPInvalidJSON(t *testing.T) {
	mockUseCase := &mockCEPInputUseCase{
		forwardCEPFunc: func(cep string) (*domain.ServiceBResponse, error) {
			return &domain.ServiceBResponse{}, nil
		},
	}

	handler := NewInputHandler(mockUseCase)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"cep":`)))
	w := httptest.NewRecorder()

	handler.ForwardCEP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `{"message":"invalid zipcode"}`, w.Body.String())
}

func TestForwardCEPInvalidZipcodeType(t *testing.T) {
	mockUseCase := &mockCEPInputUseCase{
		forwardCEPFunc: func(cep string) (*domain.ServiceBResponse, error) {
			return &domain.ServiceBResponse{}, nil
		},
	}

	handler := NewInputHandler(mockUseCase)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"cep":29902555}`)))
	w := httptest.NewRecorder()

	handler.ForwardCEP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `{"message":"invalid zipcode"}`, w.Body.String())
}

func TestForwardCEPInvalidZipcodeFormat(t *testing.T) {
	mockUseCase := &mockCEPInputUseCase{
		forwardCEPFunc: func(cep string) (*domain.ServiceBResponse, error) {
			return nil, domain.ErrInvalidCEPFormat
		},
	}

	handler := NewInputHandler(mockUseCase)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"cep":"2990255"}`)))
	w := httptest.NewRecorder()

	handler.ForwardCEP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `{"message":"invalid zipcode"}`, w.Body.String())
}

func TestForwardCEPInternalError(t *testing.T) {
	mockUseCase := &mockCEPInputUseCase{
		forwardCEPFunc: func(cep string) (*domain.ServiceBResponse, error) {
			return nil, errors.New("unexpected")
		},
	}

	handler := NewInputHandler(mockUseCase)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"cep":"29902555"}`)))
	w := httptest.NewRecorder()

	handler.ForwardCEP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message":"internal server error"}`, w.Body.String())
}
