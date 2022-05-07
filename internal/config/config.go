package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort  int
	Environment string
	JWTSecret   string
	DBPath      string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort:  serverPort,
		Environment: getEnv("ENVIRONMENT", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		DBPath:      getEnv("DB_PATH", "wisetech_lms.db"),
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
