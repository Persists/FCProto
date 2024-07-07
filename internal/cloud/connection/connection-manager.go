package connection

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Persists/fcproto/internal/cloud/config"
	"github.com/Persists/fcproto/internal/cloud/database"
	"github.com/Persists/fcproto/internal/cloud/database/models/entities"
	"github.com/Persists/fcproto/internal/cloud/messages-utils"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// ConnectionManager manages the lifecycle of a network connection
type ConnectionManager struct {
	listener     net.Listener
	db           database.DB
	connections  map[string]net.Conn // Map keys are client IP addresses
	connectionsM sync.Mutex
}

// Start and Stop methods for the ConnectionManager
func (cm *ConnectionManager) Init(config *config.ServerConfig, db database.DB) error {
	cm.connections = make(map[string]net.Conn)
	cm.db = db
	listener, err := net.Listen("tcp", config.SocketAddr)

	if err != nil {
		log.Printf("failed to listen: %v", err)
		return err
	}
	cm.listener = listener
	return nil
}

func (cm *ConnectionManager) Start() (err error) {
	messages_utils.NotifyAllClients(&cm.db)
	cm.StartSensorDataRoutine()
	for {
		conn, err := cm.listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v", err)
			return err
		}

		client, err := cm.db.InsertClient(conn.RemoteAddr().String())
		if err != nil {
			log.Printf("failed to insert client into database: %v", err)
			return err
		}

		cm.addConnection(conn)
		go func() {
			if closeErr := cm.handleConn(conn, client); closeErr != nil {
				fmt.Printf("Failed to handle connection: %v", closeErr)
				err = closeErr
			}
		}()
	}
}

func (cm *ConnectionManager) Stop() error {
	cm.connectionsM.Lock()
	defer cm.connectionsM.Unlock()

	var errList []error

	for _, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
			errList = append(errList, err)
		}
	}

	if err := cm.listener.Close(); err != nil {
		log.Printf("Failed to close listener: %v", err)
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		return fmt.Errorf("encountered errors while stopping connections: %v", errList)
	}

	return nil
}

func (cm *ConnectionManager) StartSensorDataRoutine() {
	go func() { // Added missing go routine initialization
		for {
			messages, err := cm.db.GetRecentSensorMessages()
			if err != nil {
				log.Printf("Error fetching recent sensor messages: %v", err)
				continue
			}
			fmt.Println("Starting sensor data routine")

			cm.connectionsM.Lock()
			for _, msg := range messages {
				conn, exists := cm.connections[msg.ClientIpAddr]
				if exists {
					sendMsg := models.NewMessage(models.All, models.NewSensorMessage("LOL"))
					utils.SendMessage(conn, sendMsg)
				}
			}
			cm.connectionsM.Unlock()

			// Wait for a while before checking again
			time.Sleep(10 * time.Second)
		}
	}()
}

func (cm *ConnectionManager) handleConn(conn net.Conn, client *entities.ClientEntity) (err error) {
	defer func(conn net.Conn) {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Printf("Failed to close connection: %v", closeErr)
			err = closeErr
			return
		}
		cm.removeConnection(conn)
	}(conn)

	stopTickerChan := make(chan bool)
	go utils.StartTicker(30*time.Second, func() any {
		messages_utils.UpdateLastSeen(&cm.db, client)
		return nil
	}, stopTickerChan, nil)

	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)

	for {
		var message models.Message
		err := decoder.Decode(&message)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client: %v", conn.RemoteAddr())
				stopTickerChan <- true
				close(stopTickerChan)
				return err
			}
			log.Printf("Failed to decode JSON: %v", err)
			continue
		}

		switch message.Topic {
		case models.Sensor:
			err := messages_utils.InsertSensorMessage(&cm.db, message.Payload, client)
			if err != nil {
				log.Printf("Failed to insert sensor message into database: %v", err)
				return err
			}
		case models.Heartbeat:
			err := messages_utils.UpdateNotifyAddr(&cm.db, message.Payload, client)
			if err != nil {
				log.Printf("Failed to update notify addr: %v", err)
				return err
			}
		}
	}
}

func (cm *ConnectionManager) removeConnection(conn net.Conn) {
	cm.connectionsM.Lock()
	defer cm.connectionsM.Unlock()
	delete(cm.connections, conn.RemoteAddr().String())
}

func (cm *ConnectionManager) addConnection(conn net.Conn) {
	cm.connectionsM.Lock()
	defer cm.connectionsM.Unlock()
	cm.connections[conn.RemoteAddr().String()] = conn
}
