package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xavierpms/service-a/internal/config"
	"github.com/xavierpms/service-a/internal/infra/repository"
	"github.com/xavierpms/service-a/internal/infra/validator"
	"github.com/xavierpms/service-a/internal/infra/webserver/handlers"
	"github.com/xavierpms/service-a/internal/usecase"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf(
		"Config loaded: port=%q service_b_url=%q",
		cfg.Port,
		cfg.ServiceBURL,
	)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	cepValidator := validator.NewCEPValidator()
	serviceBRepository := repository.NewServiceBRepository(cfg.ServiceBURL)
	forwardCEPUseCase := usecase.NewForwardCEPUseCase(serviceBRepository, cepValidator)
	inputHandler := handlers.NewInputHandler(forwardCEPUseCase)

	router.Post("/", inputHandler.ForwardCEP)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
