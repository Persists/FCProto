package cloud

import (
	"log"
	"time"

	"github.com/Persists/fcproto/internal/cloud/config"
	"github.com/Persists/fcproto/internal/cloud/database"
	"github.com/Persists/fcproto/internal/shared/connection"
	"github.com/Persists/fcproto/internal/shared/models"

	messages_utils "github.com/Persists/fcproto/internal/cloud/messages-utils"
)

type CloudClient struct {
	db database.DB
	lc *connection.ListenerClient
}

func NewClient() *CloudClient {
	return &CloudClient{}
}

func (cc *CloudClient) onReceive(message *models.Message, connectionClient *connection.ConnectionClient) {
	if message.Topic == models.Sensor {
		err := messages_utils.InsertSensorMessage(&cc.db, message.Payload, connectionClient.RemoteAddress())
		if err != nil {
			log.Printf("Failed to insert sensor message into database: %v", err)
		}
	}
}

func (cc *CloudClient) Init(config *config.ServerConfig) error {
	dbc := database.NewClient()
	err := dbc.Init(config)
	if err != nil {
		return err
	}

	cc.db = dbc.GetDB()
	cc.lc = connection.Listen(config.SocketAddr, cc.onReceive)

	cc.StartInformer()

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
