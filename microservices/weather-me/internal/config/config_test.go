package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvFromParentDirectory(t *testing.T) {
	originalServiceBURL, hadServiceBURL := os.LookupEnv("SERVICE_B_URL")
	originalPort, hadPort := os.LookupEnv("PORT")

	_ = os.Unsetenv("SERVICE_B_URL")
	_ = os.Unsetenv("PORT")

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
	}()

	tmpRoot := t.TempDir()
	projectRoot := filepath.Join(tmpRoot, "project")
	nestedDir := filepath.Join(projectRoot, "cmd", "server")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested directories: %v", err)
	}

	envContent := "SERVICE_B_URL=http://service-b:8080\nPORT=8081\n"
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
}

func TestLoadConfigUsesDefaultsWhenVarsAreMissingOrEmpty(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("SERVICE_B_URL", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.Port != defaultPort {
		t.Fatalf("expected default port %q, got %q", defaultPort, cfg.Port)
	}

	if cfg.ServiceBURL != defaultServiceBURL {
		t.Fatalf("expected default service-b URL %q, got %q", defaultServiceBURL, cfg.ServiceBURL)
	}
}

func TestLoadConfigPrefersEnvValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("SERVICE_B_URL", "http://localhost:8080")

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
}
