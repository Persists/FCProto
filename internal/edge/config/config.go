package config

import (
	"log"
	"strconv"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
	"github.com/joho/godotenv"
)

type ClientConfig struct {
	NotifyPort string
	SendPort   int

	*models.BaseEnv
}

// LoadConfig loads the client configuration from the environment for the edge
func LoadConfig() (*ClientConfig, error) {
	err := godotenv.Load("..edge.env")

	port, err := strconv.Atoi(utils.GetEnv("SEND_PORT", "5557"))
	if err != nil {
		log.Printf("Failed to parse SEND_PORT: %v", err)
	}

	config := &ClientConfig{
		BaseEnv:    &models.BaseEnv{SocketAddr: utils.GetEnv("SOCKET_ADDR", "localhost:5555")},
		NotifyPort: utils.GetEnv("NOTIFY_PORT", "5556"),
		SendPort:   port,
	}

	return config, nil
}
