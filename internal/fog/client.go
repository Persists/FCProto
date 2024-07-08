package fog

import (
	"fmt"
	"log"

	client_config "github.com/Persists/fcproto/internal/client/client-config"
	"github.com/Persists/fcproto/internal/shared/connection"
	"github.com/Persists/fcproto/pkg/sensors"
)

type FogClient struct {
	serverConfig *client_config.ClientConfig
	cc           *connection.ConnectionClient
	sc           *sensors.SensorClient

	stopChan chan bool
}

func NewClient() *FogClient {
	return &FogClient{
		serverConfig: &client_config.ClientConfig{},
		sc:           sensors.NewClient(),
	}
}

func (fc *FogClient) Init() error {
	config, err := client_config.LoadConfig()

	if err != nil {
		log.Printf("failed to load config: %v", err)
		return err
	}

	client := connection.Connect(config.SocketAddr)
	fc.cc = client

	fc.sc.Start(fc.stopChan, fc.cc.Send)

	for {
		serverMessage := fc.cc.Receive()

		fmt.Printf("Received message from server: Topic: %s, Payload: %s\n", serverMessage.Topic, serverMessage.Payload)
	}

	return nil
}
