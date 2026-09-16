package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	DBSSLMode         string
	JWTSecret         string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
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
		JWTSecret: getEnv(
			"JWT_SECRET",
			"vigipeconha-development-secret",
		),
		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2BucketName:      os.Getenv("R2_BUCKET_NAME"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
