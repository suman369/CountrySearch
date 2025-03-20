package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port                 string
	Timeout              int
	RestCountriesBaseURL string
}

// LoadConfig loads the configuration from environment variables or uses defaults
func LoadConfig() *Config {
	port := getEnv("PORT", "8000")
	timeout, _ := strconv.Atoi(getEnv("TIMEOUT", "10"))
	restCountriesBaseURL := getEnv("REST_COUNTRIES_URL", "https://restcountries.com/v3.1")

	return &Config{
		Port:                 port,
		Timeout:              timeout,
		RestCountriesBaseURL: restCountriesBaseURL,
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
