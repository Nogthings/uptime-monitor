package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
	JWTSecret     string
}

// Load reads configuration from environment variables and returns a Config struct.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}

	return cfg, nil
}
