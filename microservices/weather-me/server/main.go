package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xavierpms/service-a/internal/config"
	"github.com/xavierpms/service-a/internal/infra/repository"
	"github.com/xavierpms/service-a/internal/infra/validator"
	"github.com/xavierpms/service-a/internal/infra/webserver/handlers"
	"github.com/xavierpms/service-a/internal/observability"
	"github.com/xavierpms/service-a/internal/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf(
		"Config loaded: port=%q service_b_url=%q zipkin_endpoint=%q",
		cfg.Port,
		cfg.ServiceBURL,
		cfg.ZipkinEndpoint,
	)

	tracerProvider, err := observability.InitTracerProvider(context.Background(), "weather-me", cfg.ZipkinEndpoint)
	if err != nil {
		log.Fatalf("failed to initialize tracer provider: %v", err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	cepValidator := validator.NewCEPValidator()
	serviceBRepository := repository.NewServiceBRepository(cfg.ServiceBURL)
	forwardCEPUseCase := usecase.NewForwardCEPUseCase(serviceBRepository, cepValidator)
	inputHandler := handlers.NewInputHandler(forwardCEPUseCase)

	router.Post("/", inputHandler.ForwardCEP)

	log.Printf("Starting server on port %s", cfg.Port)
	handler := otelhttp.NewHandler(router, "weather-me.http.server")
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
