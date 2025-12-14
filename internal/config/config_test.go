package config

import (
	"os"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Unset env vars to ensure we are testing default values
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DB_PATH")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check default values
	if cfg.ServerPort != 8080 {
		t.Errorf("Expected ServerPort to be 8080, got %d", cfg.ServerPort)
	}
	if cfg.Environment != "development" {
		t.Errorf("Expected Environment to be 'development', got %s", cfg.Environment)
	}
	if cfg.JWTSecret != "your-secret-key" {
		t.Errorf("Expected JWTSecret to be 'your-secret-key', got %s", cfg.JWTSecret)
	}
	if cfg.DBPath != "wisetech_lms.db" {
		t.Errorf("Expected DBPath to be 'wisetech_lms.db', got %s", cfg.DBPath)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("JWT_SECRET", "a-different-secret")
	os.Setenv("DB_PATH", "/tmp/test.db")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check values from env vars
	if cfg.ServerPort != 9090 {
		t.Errorf("Expected ServerPort to be 9090, got %d", cfg.ServerPort)
	}
	if cfg.Environment != "production" {
		t.Errorf("Expected Environment to be 'production', got %s", cfg.Environment)
	}
	if cfg.JWTSecret != "a-different-secret" {
		t.Errorf("Expected JWTSecret to be 'a-different-secret', got %s", cfg.JWTSecret)
	}
	if cfg.DBPath != "/tmp/test.db" {
		t.Errorf("Expected DBPath to be '/tmp/test.db', got %s", cfg.DBPath)
	}

	// Unset env vars
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DB_PATH")
}

func TestLoadConfig_FromDotEnv(t *testing.T) {
	// Create a temporary .env file
	content := []byte("SERVER_PORT=7070\nENVIRONMENT=staging\n")
	err := os.WriteFile(".env", content, 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check values from .env file
	if cfg.ServerPort != 7070 {
		t.Errorf("Expected ServerPort to be 7070, got %d", cfg.ServerPort)
	}
	if cfg.Environment != "staging" {
		t.Errorf("Expected Environment to be 'staging', got %s", cfg.Environment)
	}
}
