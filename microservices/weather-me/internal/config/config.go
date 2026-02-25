package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	ServiceBURL string
}

const (
	defaultPort        = "8081"
	defaultServiceBURL = "http://localhost:8080"
)

// LoadConfig loads environment variables and returns a Config struct.
func LoadConfig() (*Config, error) {
	loadDotEnv()

	return &Config{
		Port:        getEnv("PORT", defaultPort),
		ServiceBURL: getEnv("SERVICE_B_URL", defaultServiceBURL),
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
