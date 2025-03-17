package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	BaseURL         string
	RateLimitPerMin int
	AllowedOrigins  []string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		Port:            getEnvOrDefault("PORT", ":8080"),
		BaseURL:         getEnvOrDefault("BASE_URL", "https://baak.gunadarma.ac.id"),
		RateLimitPerMin: getEnvIntOrDefault("RATE_LIMIT_PER_MIN", 60),
		AllowedOrigins:  getEnvSliceOrDefault("ALLOWED_ORIGINS", []string{"*"}),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return []string{value}
	}
	return defaultValue
}
