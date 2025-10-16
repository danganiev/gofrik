package api

import (
	"os"
)

// Config holds all server configuration
type Config struct {
	Port       string
	DatabaseURL string
	Host       string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/gofrik?sslmode=disable"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = ""
	}

	return &Config{
		Port:       port,
		DatabaseURL: dbURL,
		Host:       host,
	}
}

