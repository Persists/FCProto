package config

import (
	database_models "github.com/Persists/fcproto/internal/cloud/database/models"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
	"github.com/joho/godotenv"
	"log"
)

type ServerConfig struct {
	*models.BaseEnv
	*database_models.PostgresEnv
}

func LoadConfig() (*ServerConfig, error) {
	err := godotenv.Load(".env", ".cloud.env")

	if err != nil {
		log.Printf("No .cloud.env file found, using fallback variables: %v\n", err)
	}

	baseEnv := &models.BaseEnv{
		SocketAddr: utils.GetEnv("SOCKET_ADDR", "localhost:5555"),
	}

	postgresEnv := &database_models.PostgresEnv{
		Addr:     utils.GetEnv("POSTGRES_ADDR", "localhost:5432"),
		User:     utils.GetEnv("POSTGRES_USER", "fog"),
		Password: utils.GetEnv("POSTGRES_PASSWORD", "fog"),
		Database: utils.GetEnv("POSTGRES_DATABASE", "fog"),
	}

	config := &ServerConfig{
		baseEnv,
		postgresEnv,
	}

	return config, nil
}
