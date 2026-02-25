package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xavierpms/service-a/internal/domain"
)

// InputHandler handles CEP input requests.
type InputHandler struct {
	useCase domain.CEPInputUseCase
}

// CEPInputRequest represents the expected POST body.
type CEPInputRequest struct {
	CEP interface{} `json:"cep"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Message string `json:"message"`
}

// NewInputHandler creates a new input handler.
func NewInputHandler(useCase domain.CEPInputUseCase) *InputHandler {
	return &InputHandler{useCase: useCase}
}

// ForwardCEP handles POST / requests and forwards CEP to weather-by-city.
func (h *InputHandler) ForwardCEP(w http.ResponseWriter, r *http.Request) {
	var request CEPInputRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeInvalidZipcode(w)
		return
	}

	cep, ok := request.CEP.(string)
	if !ok {
		h.writeInvalidZipcode(w)
		return
	}

	response, err := h.useCase.ForwardCEP(r.Context(), cep)
	if err != nil {
		h.handleError(w, err)
		return
	}

	contentType := response.ContentType
	if contentType == "" {
		contentType = "application/json"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(response.Body)
}

func (h *InputHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInvalidCEPFormat:
		h.writeInvalidZipcode(w)
	default:
		log.Printf("error forwarding CEP: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "internal server error"})
	}
}

func (h *InputHandler) writeInvalidZipcode(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
}
