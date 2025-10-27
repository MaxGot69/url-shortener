package config

import (
	"os"
)

type Config struct {
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
	JWTSecret        string
	LogLevel         string
}

func LoadConfig() Config {
	return Config{
		Port:          getEnv("PORT", "8081"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "maxim"),
		DBPassword:    getEnv("DB_PASSWORD", "secret"),
		DBName:        getEnv("DB_NAME", "urlshortener"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
