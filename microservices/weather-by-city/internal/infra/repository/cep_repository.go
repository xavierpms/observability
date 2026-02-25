package repository

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/xavierpms/weather-by-city/internal/domain"
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
}

// NewCEPRepository creates a new CEP repository
func NewCEPRepository(apiURL string) domain.CEPRepository {
	return &CEPRepositoryImpl{
		apiURL: apiURL,
	}
}

// GetCEPData fetches the data for a given CEP
func (r *CEPRepositoryImpl) GetCEPData(cep string) (*domain.CEPData, error) {
	// Build the URL
	requestURL := r.apiURL + "/" + cep + "/json/"
	log.Printf("calling ViaCEP API: url=%s cep=%s", requestURL, cep)

	// Make the request
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("ViaCEP API request error: cep=%s err=%v", cep, err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("ViaCEP API response: cep=%s status=%d", cep, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ViaCEP API read response error: cep=%s err=%v", cep, err)
		return nil, err
	}

	// Unmarshal the response
	var viaCepData ViaCEPResponse
	err = json.Unmarshal(body, &viaCepData)
	if err != nil {
		log.Printf("ViaCEP API parse response error: cep=%s err=%v", cep, err)
		return nil, err
	}

	// Validate if the CEP was found
	if viaCepData.Erro {
		log.Printf("ViaCEP API returned not found: cep=%s", cep)
		return nil, errors.New("CEP not found in ViaCEP")
	}
	log.Printf("ViaCEP API request succeeded: cep=%s city=%s", cep, viaCepData.Localidade)

	return &domain.CEPData{
		CEP:    viaCepData.CEP,
		City:   viaCepData.Localidade,
		Region: viaCepData.UF,
	}, nil
}
