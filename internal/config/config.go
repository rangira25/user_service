package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimeZone string
	Port       string
}

// LoadConfig loads environment variables into a Config struct
func LoadConfig() *Config {
	_ = godotenv.Load() // load from .env if present

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBTimeZone: os.Getenv("DB_TIMEZONE"),
		Port:       os.Getenv("PORT"),
	}

	// Basic check
	if cfg.DBHost == "" || cfg.DBUser == "" {
		log.Fatal(" Missing database environment variables")
	}

	return cfg
}
