package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabasePath    string
	Port            string
	LogLevel        string
	AllowedOrigins  string
	DefaultLanguage string
}

func Load() *Config {
	return &Config{
		DatabasePath:    getEnv("DATABASE_PATH", "./hadith.db"),
		Port:            getEnv("PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		AllowedOrigins:  getEnv("ALLOWED_ORIGINS", "*"),
		DefaultLanguage: getEnv("DEFAULT_LANGUAGE", "ar"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
