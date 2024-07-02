package server_config

import (
	database_models "github.com/Persists/fcproto/internal/server/database/models"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ServerConfig struct {
	*models.BaseEnv
	*database_models.PostgresEnv
}

func LoadConfig() (*ServerConfig, error) {
	err := godotenv.Load(".env", ".server.env")

	if err != nil {
		log.Printf("No ..server.env file found, using fallback variables: %v\n", err)
	}

	baseEnv := &models.BaseEnv{
		SocketAddr: getEnv("SOCKET_ADDR", "localhost:5555"),
	}

	postgresEnv := &database_models.PostgresEnv{
		User:     getEnv("POSTGRES_USER", "fog"),
		Password: getEnv("POSTGRES_PASSWORD", "fog"),
		Database: getEnv("POSTGRES_DATABASE", "fog"),
	}

	config := &ServerConfig{
		baseEnv,
		postgresEnv,
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
