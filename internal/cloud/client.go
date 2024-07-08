package cloud

import (
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
		log.Printf("Failed to start database client: %v", err)
	}

	db := cc.dbc.GetDB()
	recentMsg, _ := database.GetRecentSensorMessages(&db)
	sensors.SortMessageData(recentMsg)

	cc.lc = connection.Listen(cc.socketAddress, cc.onReceive)

	cc.StartInformer()

	return nil

}

func (cc *CloudClient) Init(config *config.ServerConfig) error {
	cc.dbc = database.NewClient()
	cc.socketAddress = config.SocketAddr

	err := cc.dbc.Init(config)
	if err != nil {
		return err
	}

	return nil
}

func (cc *CloudClient) StartInformer() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			cc.InformFog()
		}
	}()

}

func (cc *CloudClient) InformFog() {
	message := models.Message{
		Topic:   models.Fog,
		Payload: &map[string]interface{}{"message": "Hello Fog!"},
	}
	cc.lc.BroadcastMsg(message)
}
