package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config represents the application configuration.
type Config struct {
	DatabaseURL    string
	ExternalAPIURL string
}

// @Summary Load application configuration
// @Description Loads application configuration from environment variables
// @Tags config
// @Produce json
// @Success 200 {object} Config
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		ExternalAPIURL: os.Getenv("EXTERNAL_API_URL"),
	}

	return config
}
