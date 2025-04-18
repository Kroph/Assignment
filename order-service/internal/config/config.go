package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Server struct {
		Port string
	}
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{}

	config.DB.Host = getEnv("ORDERS_DB_HOST", "localhost")
	config.DB.Port = getEnv("ORDERS_DB_PORT", "5432")
	config.DB.User = getEnv("ORDERS_DB_USER", "postgres")
	config.DB.Password = getEnv("ORDERS_DB_PASSWORD", "postgres")
	config.DB.Name = getEnv("ORDERS_DB_NAME", "orders")

	config.Server.Port = getEnv("ORDERS_PORT", "8081")

	return config
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
