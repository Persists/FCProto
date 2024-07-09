package cloud

import (
	"fmt"
	"github.com/Persists/fcproto/internal/shared/utils"
	"log"
	"time"

	"github.com/Persists/fcproto/internal/cloud/config"
	"github.com/Persists/fcproto/internal/cloud/database"
	"github.com/Persists/fcproto/internal/shared/connection"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/pkg/sensors"
)

type CloudClient struct {
	dbc *database.DBClient
	lc  *connection.ListenerClient

	socketAddress string
}

func NewClient() *CloudClient {
	return &CloudClient{}
}

func (cc *CloudClient) onReceive(message *models.Message, connectionClient *connection.ConnectionClient) {
	db := cc.dbc.GetDB()

	if message.Topic == models.Sensor {
		err := database.InsertSensorMessage(&db, message.Payload, connectionClient.RemoteAddress())
		if err != nil {
			log.Printf("Failed to insert sensor message into database: %v", err)
		}
	}
}

func (cc *CloudClient) Start() error {
	err := cc.dbc.Start()
	if err != nil {
		message := fmt.Sprintf("Failed to start database client: %v", err)
		log.Println(utils.Colorize(utils.Yellow, message))
	}

	cc.lc = connection.Listen(cc.socketAddress, cc.onReceive)

	cc.StartInformer()

	return nil

}

func (cc *CloudClient) Init(config *config.ServerConfig) error {
	cc.dbc = database.NewClient()
	cc.socketAddress = config.SocketAddr

	cc.dbc.Init(config)

	return nil
}

func (cc *CloudClient) StartInformer() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			cc.InformFog()
		}
	}()

}

func (cc *CloudClient) InformFog() {
	db := cc.dbc.GetDB()
	for ip, conn := range cc.lc.Connections {
		sensorData, err := database.GetRecentSensorMessages(&db, ip)

		if err != nil {
			log.Printf("Failed to get recent sensor messages: %v", err)
			continue
		}

		VSD, MSD := sensors.MapMessageData(sensorData)

		mappedVSD := make([]float64, len(VSD))
		for i, data := range VSD {
			mappedVSD[i] = data.Temperature
		}

		mappedMSD := make([]float64, len(MSD))
		for i, data := range MSD {
			mappedMSD[i] = float64(data.Total)
		}

		avgVSD, devVSD, meanVSD, minVSD, maxVSD := sensors.Analyse(mappedVSD)
		avgMSD, devMSD, meanMSD, minMSD, maxMSD := sensors.Analyse(mappedMSD)

		payload := models.AnalysisPayload{
			VirtualSensorData: models.AnalysisVirtualSensorData{
				Average:   avgVSD,
				Deviation: devVSD,
				Mean:      meanVSD,
				Min:       minVSD,
				Max:       maxVSD,
			},
			MemorySensorData: models.AnalysisMemorySensorData{
				Average:   avgMSD,
				Deviation: devMSD,
				Mean:      meanMSD,
				Min:       minMSD,
				Max:       maxMSD,
			},
		}

		msg := models.NewMessage(models.Analysis, payload)

		conn.Send(msg)

	}
}
