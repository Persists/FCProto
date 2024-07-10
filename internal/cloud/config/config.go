package config

import (
	database_models "github.com/Persists/fcproto/internal/cloud/database/models"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
	"github.com/joho/godotenv"
)

// ServerConfig represents the server configuration
type ServerConfig struct {
	*models.BaseEnv
	*database_models.PostgresEnv
}

// LoadConfig loads the server configuration from the environment
func LoadConfig() (*ServerConfig, error) {
	_ = godotenv.Load(".cloud.env")

	baseEnv := &models.BaseEnv{
		SocketAddr: utils.GetEnv("SOCKET_ADDR", ":5555"),
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
