package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Services struct {
		Inventory struct {
			URL string
		}
		Order struct {
			URL string
		}
	}
	Auth struct {
		Secret        string
		ExpiryMinutes int
	}
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{}

	config.Server.Port = getEnv("GATEWAY_PORT", "8000")

	config.Services.Inventory.URL = getEnv("INVENTORY_SERVICE_URL", "http://localhost:8080/api/v1")
	config.Services.Order.URL = getEnv("ORDER_SERVICE_URL", "http://localhost:8081/api/v1")

	config.Auth.Secret = getEnv("AUTH_SECRET", "your-default-secret-key-change-in-production")
	expiryStr := getEnv("AUTH_EXPIRY_MINUTES", "60")
	expiryMinutes, err := strconv.Atoi(expiryStr)
	if err != nil {
		expiryMinutes = 60
	}
	config.Auth.ExpiryMinutes = expiryMinutes

	return config
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
