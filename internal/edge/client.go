package edge

import (
	"fmt"
	"log"

	"github.com/Persists/fcproto/internal/edge/config"
	"github.com/Persists/fcproto/internal/shared/utils"

	"github.com/Persists/fcproto/internal/shared/connection"
	"github.com/Persists/fcproto/pkg/sensors"
)

type EdgeClient struct {
	serverConfig *config.ClientConfig
	cc           *connection.ConnectionClient
	sc           *sensors.SensorClient

	stopChan chan bool
}

func NewClient() *EdgeClient {
	return &EdgeClient{
		serverConfig: &config.ClientConfig{},
		sc:           sensors.NewClient(),
	}
}

// Start initializes the fog client
func (ec *EdgeClient) Start() error {
	config, err := config.LoadConfig()

	if err != nil {
		log.Printf("failed to load config: %v", err)
		return err
	}
	client := connection.Connect(config.SocketAddr)
	ec.cc = client

	ec.sc.Start(ec.stopChan, ec.cc.Send)

	fmt.Println("Fog client initialized")
	for {
		serverMessage := ec.cc.Receive()
		formattedServerMessage := utils.FormatCloudAnalysisData(serverMessage)
		log.Println(utils.Colorize(utils.Yellow, formattedServerMessage))
	}
}
