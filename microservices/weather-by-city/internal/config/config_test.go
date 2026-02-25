package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvFromParentDirectory(t *testing.T) {
	originalWeatherAPIKey, hadWeatherAPIKey := os.LookupEnv("WEATHER_API_KEY")
	originalPort, hadPort := os.LookupEnv("PORT")

	_ = os.Unsetenv("WEATHER_API_KEY")
	_ = os.Unsetenv("PORT")

	defer func() {
		if hadWeatherAPIKey {
			_ = os.Setenv("WEATHER_API_KEY", originalWeatherAPIKey)
		} else {
			_ = os.Unsetenv("WEATHER_API_KEY")
		}

		if hadPort {
			_ = os.Setenv("PORT", originalPort)
		} else {
			_ = os.Unsetenv("PORT")
		}
	}()

	tmpRoot := t.TempDir()
	projectRoot := filepath.Join(tmpRoot, "project")
	nestedDir := filepath.Join(projectRoot, "cmd", "server")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested directories: %v", err)
	}

	envContent := "WEATHER_API_KEY=test_parent_key\nPORT=9999\n"
	if err := os.WriteFile(filepath.Join(projectRoot, ".env"), []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	if err := os.Chdir(nestedDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWd)
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.WeatherAPIKey != "test_parent_key" {
		t.Fatalf("expected WeatherAPIKey to be loaded from parent .env, got %q", cfg.WeatherAPIKey)
	}

	if cfg.Port != "9999" {
		t.Fatalf("expected Port to be loaded from parent .env, got %q", cfg.Port)
	}
}

func TestLoadConfigUsesDefaultWeatherAPIKeyWhenMissingOrEmpty(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.WeatherAPIKey != defaultWeatherAPIKey {
		t.Fatalf("expected default WeatherAPIKey %q, got %q", defaultWeatherAPIKey, cfg.WeatherAPIKey)
	}
}

func TestLoadConfigPrefersWeatherAPIKeyEnv(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "env_secret_key")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.WeatherAPIKey != "env_secret_key" {
		t.Fatalf("expected WeatherAPIKey loaded from env var, got %q", cfg.WeatherAPIKey)
	}
}

func TestLoadConfigUsesDefaultsWhenVarsAreMissingOrEmpty(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("WEATHER_API_URL", "")
	t.Setenv("VIA_CEP_URL", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.Port != defaultPort {
		t.Fatalf("expected default port %q, got %q", defaultPort, cfg.Port)
	}

	if cfg.WeatherAPIURL != defaultWeatherAPIURL {
		t.Fatalf("expected default weather API URL %q, got %q", defaultWeatherAPIURL, cfg.WeatherAPIURL)
	}

	if cfg.ViaCEPURL != defaultViaCEPURL {
		t.Fatalf("expected default ViaCEP URL %q, got %q", defaultViaCEPURL, cfg.ViaCEPURL)
	}
}
