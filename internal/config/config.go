package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DBConfig
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DB: DBConfig{
			User:     mustGet("DB_USER"),
			Password: mustGet("DB_PASSWORD"),
			Host:     mustGet("DB_HOST"),
			Port:     mustGet("DB_PORT"),
			Name:     mustGet("DB_NAME"),
			SSLMode:  mustGet("DB_SSLMODE"),
		},
	}
	return cfg
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}
