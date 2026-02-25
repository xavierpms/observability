package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvFromParentDirectory(t *testing.T) {
	originalServiceBURL, hadServiceBURL := os.LookupEnv("SERVICE_B_URL")
	originalPort, hadPort := os.LookupEnv("PORT")
	originalZipkinEndpoint, hadZipkinEndpoint := os.LookupEnv("ZIPKIN_ENDPOINT")

	_ = os.Unsetenv("SERVICE_B_URL")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("ZIPKIN_ENDPOINT")

	defer func() {
		if hadServiceBURL {
			_ = os.Setenv("SERVICE_B_URL", originalServiceBURL)
		} else {
			_ = os.Unsetenv("SERVICE_B_URL")
		}

		if hadPort {
			_ = os.Setenv("PORT", originalPort)
		} else {
			_ = os.Unsetenv("PORT")
		}

		if hadZipkinEndpoint {
			_ = os.Setenv("ZIPKIN_ENDPOINT", originalZipkinEndpoint)
		} else {
			_ = os.Unsetenv("ZIPKIN_ENDPOINT")
		}
	}()

	tmpRoot := t.TempDir()
	projectRoot := filepath.Join(tmpRoot, "project")
	nestedDir := filepath.Join(projectRoot, "cmd", "server")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested directories: %v", err)
	}

	envContent := "SERVICE_B_URL=http://service-b:8080\nPORT=8081\nZIPKIN_ENDPOINT=http://zipkin:9411/api/v2/spans\n"
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

	if cfg.ServiceBURL != "http://service-b:8080" {
		t.Fatalf("expected ServiceBURL loaded from parent .env, got %q", cfg.ServiceBURL)
	}

	if cfg.Port != "8081" {
		t.Fatalf("expected Port loaded from parent .env, got %q", cfg.Port)
	}

	if cfg.ZipkinEndpoint != "http://zipkin:9411/api/v2/spans" {
		t.Fatalf("expected ZipkinEndpoint loaded from parent .env, got %q", cfg.ZipkinEndpoint)
	}
}

func TestLoadConfigUsesDefaultsWhenVarsAreMissingOrEmpty(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("SERVICE_B_URL", "")
	t.Setenv("ZIPKIN_ENDPOINT", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.Port != defaultPort {
		t.Fatalf("expected default port %q, got %q", defaultPort, cfg.Port)
	}

	if cfg.ServiceBURL != defaultServiceBURL {
		t.Fatalf("Expected default weather-by-city URL %q, got %q", defaultServiceBURL, cfg.ServiceBURL)
	}

	if cfg.ZipkinEndpoint != defaultZipkinEndpoint {
		t.Fatalf("Expected default zipkin endpoint %q, got %q", defaultZipkinEndpoint, cfg.ZipkinEndpoint)
	}
}

func TestLoadConfigPrefersEnvValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("SERVICE_B_URL", "http://localhost:8080")
	t.Setenv("ZIPKIN_ENDPOINT", "http://localhost:9411/api/v2/spans")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.Port != "9090" {
		t.Fatalf("expected Port loaded from env var, got %q", cfg.Port)
	}

	if cfg.ServiceBURL != "http://localhost:8080" {
		t.Fatalf("expected ServiceBURL loaded from env var, got %q", cfg.ServiceBURL)
	}

	if cfg.ZipkinEndpoint != "http://localhost:9411/api/v2/spans" {
		t.Fatalf("expected ZipkinEndpoint loaded from env var, got %q", cfg.ZipkinEndpoint)
	}
}
