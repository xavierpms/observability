package repository

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/xavierpms/service-a/internal/domain"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ServiceBRepositoryImpl implements domain.ServiceBRepository.
type ServiceBRepositoryImpl struct {
	apiURL string
	client *http.Client
}

// NewServiceBRepository creates a new repository for weather-by-city HTTP calls.
func NewServiceBRepository(apiURL string) domain.ServiceBRepository {
	return &ServiceBRepositoryImpl{
		apiURL: strings.TrimRight(apiURL, "/"),
		client: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

// ForwardCEP forwards CEP to weather-by-city and returns raw response details.
func (r *ServiceBRepositoryImpl) ForwardCEP(ctx context.Context, cep string) (*domain.ServiceBResponse, error) {
	tracer := otel.Tracer("weather-me.repository.weather-by-city")
	ctx, span := tracer.Start(ctx, "weather-by-city.request")
	defer span.End()

	requestURL := r.apiURL + "/" + cep
	span.SetAttributes(attribute.String("http.url", requestURL), attribute.String("zipcode", cep))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to build weather-by-city request")
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "weather-by-city request failed")
		return nil, err
	}
	defer resp.Body.Close()
	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to read weather-by-city response")
		return nil, err
	}

	return &domain.ServiceBResponse{
		StatusCode:  resp.StatusCode,
		Body:        body,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
}
