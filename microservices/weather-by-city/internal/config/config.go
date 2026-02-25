package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	WeatherAPIKey string
	WeatherAPIURL string
	ViaCEPURL     string
}

const (
	defaultPort          = "8080"
	// The correct place to storage the API key would be in a secret manager or environment variable, but for the sake of this exercise, we will use a default value.
	defaultWeatherAPIKey  = "d37a29f8938b4c07adb173527262402"
	defaultWeatherAPIURL = "https://api.weatherapi.com/v1/current.json"
	defaultViaCEPURL     = "https://viacep.com.br/ws"
)

// LoadConfig loads the environment variables and returns a Config struct
func LoadConfig() (*Config, error) {
	loadDotEnv()

	weatherAPIKey := strings.TrimSpace(getEnv("WEATHER_API_KEY", ""))
	if weatherAPIKey == "" {
		weatherAPIKey = defaultWeatherAPIKey
	}

	return &Config{
		Port:          getEnv("PORT", defaultPort),
		WeatherAPIKey: weatherAPIKey,
		WeatherAPIURL: getEnv("WEATHER_API_URL", defaultWeatherAPIURL),
		ViaCEPURL:     getEnv("VIA_CEP_URL", defaultViaCEPURL),
	}, nil
}

func loadDotEnv() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(wd, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			return
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			return
		}

		wd = parent
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		trimmedValue := strings.TrimSpace(value)
		if trimmedValue != "" {
			return trimmedValue
		}

		return defaultVal
	}
	return defaultVal
}
