package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	port      string
	dbURL     string
	jwtSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	config := &Config{
		port:      getEnv("PORT", "8080"),
		dbURL:     getEnv("DATABASE_URL", ""),
		jwtSecret: getEnv("JWT_SECRET", ""),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
