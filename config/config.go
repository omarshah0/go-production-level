package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	RedisURL    string
	JWTSecret   string
	ServerPort  string
	Environment string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	config := &Config{
		DatabaseUrl: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Print all config values
	fmt.Printf("Database URL: %s\n", config.DatabaseUrl)
	fmt.Printf("Redis URL: %s\n", config.RedisURL)
	fmt.Printf("JWT Secret: %s\n", config.JWTSecret)
	fmt.Printf("Server Port: %s\n", config.ServerPort)
	fmt.Printf("Environment: %s\n", config.Environment)

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
