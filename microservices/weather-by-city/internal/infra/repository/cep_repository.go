package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/xavierpms/weather-by-city/internal/domain"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ViaCEPResponse represents the response from the ViaCEP API
type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}

// CEPRepositoryImpl implement domain.CEPRepository
type CEPRepositoryImpl struct {
	apiURL string
	client *http.Client
}

// NewCEPRepository creates a new CEP repository
func NewCEPRepository(apiURL string) domain.CEPRepository {
	return &CEPRepositoryImpl{
		apiURL: apiURL,
		client: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

// GetCEPData fetches the data for a given CEP
func (r *CEPRepositoryImpl) GetCEPData(ctx context.Context, cep string) (*domain.CEPData, error) {
	tracer := otel.Tracer("weather-by-city.repository.viacep")
	ctx, span := tracer.Start(ctx, "viacep.lookup")
	defer span.End()

	// Build the URL
	requestURL := r.apiURL + "/" + cep + "/json/"
	span.SetAttributes(attribute.String("http.url", requestURL), attribute.String("zipcode", cep))
	log.Printf("calling ViaCEP API: url=%s cep=%s", requestURL, cep)

	// Make the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to build ViaCEP request")
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		log.Printf("ViaCEP API request error: cep=%s err=%v", cep, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "ViaCEP request failed")
		return nil, err
	}
	defer resp.Body.Close()
	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	log.Printf("ViaCEP API response: cep=%s status=%d", cep, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ViaCEP API read response error: cep=%s err=%v", cep, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed reading ViaCEP response")
		return nil, err
	}

	// Unmarshal the response
	var viaCepData ViaCEPResponse
	err = json.Unmarshal(body, &viaCepData)
	if err != nil {
		log.Printf("ViaCEP API parse response error: cep=%s err=%v", cep, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed parsing ViaCEP response")
		return nil, err
	}

	// Validate if the CEP was found
	if viaCepData.Erro {
		log.Printf("ViaCEP API returned not found: cep=%s", cep)
		err = errors.New("CEP not found in ViaCEP")
		span.RecordError(err)
		span.SetStatus(codes.Error, "CEP not found")
		return nil, err
	}
	log.Printf("ViaCEP API request succeeded: cep=%s city=%s", cep, viaCepData.Localidade)

	return &domain.CEPData{
		CEP:    viaCepData.CEP,
		City:   viaCepData.Localidade,
		Region: viaCepData.UF,
	}, nil
}
