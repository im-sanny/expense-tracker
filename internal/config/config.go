package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration for the application
type Config struct {
	DB DBConfig
}

// DBConfig holds the database specific configuration
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

// DSN returns the Data Source Name connection string
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

// Load reads environment variables and returns a populated Config.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

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
