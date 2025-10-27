package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT       string
	DB_URL     string
	JWT_SECRET string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	config := &Config{
		PORT:       getEnv("PORT", "8080"),
		DB_URL:     getEnv("DATABASE_URL", ""),
		JWT_SECRET: getEnv("JWT_SECRET", ""),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
