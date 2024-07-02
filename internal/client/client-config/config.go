package config

import (
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadConfig() (*models.BaseEnv, error) {
	err := godotenv.Load()

	if err != nil {
		log.Printf("No .env file found, using fallback variables: %v\n", err)
	}

	config := &models.BaseEnv{
		SocketAddr: getEnv("SOCKET_ADDR", "localhost:5555"),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
