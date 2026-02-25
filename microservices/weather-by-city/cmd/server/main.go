package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xavierpms/weather-by-city/internal/config"
	"github.com/xavierpms/weather-by-city/internal/infra/repository"
	"github.com/xavierpms/weather-by-city/internal/infra/validator"
	"github.com/xavierpms/weather-by-city/internal/infra/webserver/handlers"
	"github.com/xavierpms/weather-by-city/internal/usecase"
)

func main() {
	// Load the configurations from the .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf(
		"Config loaded: port=%q weather_api_key_set=%t weather_api_url=%q via_cep_url=%q",
		cfg.Port,
		cfg.WeatherAPIKey != "",
		cfg.WeatherAPIURL,
		cfg.ViaCEPURL,
	)
	if cfg.WeatherAPIKey == "" {
		log.Printf("WARNING: WEATHER_API_KEY is empty")
	}
	if cfg.WeatherAPIURL == "" {
		log.Printf("WARNING: WEATHER_API_URL is empty")
	}
	if cfg.ViaCEPURL == "" {
		log.Printf("WARNING: VIA_CEP_URL is empty")
	}
	if cfg.Port == "" {
		log.Printf("WARNING: PORT is empty")
	}

	// Initialize the router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Inject dependencies
	cepValidator := validator.NewCEPValidator()
	cepRepository := repository.NewCEPRepository(cfg.ViaCEPURL)
	tempRepository := repository.NewTemperatureRepository(cfg.WeatherAPIURL, cfg.WeatherAPIKey)
	getTempUseCase := usecase.NewGetTemperatureByCEP(cepRepository, tempRepository, cepValidator)
	temperatureHandler := handlers.NewTemperatureHandler(getTempUseCase)

	// Define the routes
	router.Get("/{cep}", temperatureHandler.GetTemperatureByCEP)

	// Start the server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
