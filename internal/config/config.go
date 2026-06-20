package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
	JWTSecret string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "admin"),
		DBPass:    getEnv("DB_PASSWORD", "admin123"),
		DBName:    getEnv("DB_NAME", "vigipeconha"),
		DBSSLMode: getEnv("DB_SSL_MODE", "disable"),
		JWTSecret: getEnv("JWT_SECRET", "vigipeconha-secret"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
