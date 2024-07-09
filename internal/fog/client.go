package fog

import (
	"fmt"
	"github.com/Persists/fcproto/internal/fog/config"
	"github.com/Persists/fcproto/internal/shared/utils"
	"log"

	"github.com/Persists/fcproto/internal/shared/connection"
	"github.com/Persists/fcproto/pkg/sensors"
)

type FogClient struct {
	serverConfig *config.ClientConfig
	cc           *connection.ConnectionClient
	sc           *sensors.SensorClient

	stopChan chan bool
}

func NewClient() *FogClient {
	return &FogClient{
		serverConfig: &config.ClientConfig{},
		sc:           sensors.NewClient(),
	}
}

func (fc *FogClient) Init() error {
	config, err := config.LoadConfig()

	if err != nil {
		log.Printf("failed to load config: %v", err)
		return err
	}
	client := connection.Connect(config.SocketAddr)
	fc.cc = client

	fc.sc.Start(fc.stopChan, fc.cc.Send)

	fmt.Println("Fog client initialized")
	for {
		serverMessage := fc.cc.Receive()
		formattedServerMessage := utils.FormatCloudAnalysisData(serverMessage)
		log.Println(utils.Colorize(utils.Yellow, formattedServerMessage))
	}

	return nil
}
