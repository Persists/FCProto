package connection

import (
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
	connections  map[net.Conn]struct{}
	connectionsM sync.Mutex
}

// Start and Stop methods for the ConnectionManager
func (cm *ConnectionManager) Init(config *config.ServerConfig, db database.DB) error {
	cm.connections = make(map[net.Conn]struct{})
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
				fmt.Printf("Failed to handle connection: %v", err)
				err = closeErr
			}
		}()
	}
}

func (cm *ConnectionManager) Stop() error {
	cm.connectionsM.Lock()
	defer cm.connectionsM.Unlock()

	var errList []error

	for conn := range cm.connections {
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

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client: %v", conn.RemoteAddr())
				stopTickerChan <- true
				close(stopTickerChan)
				return err
			}
			log.Printf("Failed to read from connection: %v", err)
			stopTickerChan <- true
			close(stopTickerChan)
			return err
		}

		if err != nil {
			log.Printf("Failed to insert data into database: %v", err)
			return err
		}

		var message models.Message
		fmt.Println("Received data: ", n)
		err = json.Unmarshal(buf[:n], &message)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
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
	delete(cm.connections, conn)
}

func (cm *ConnectionManager) addConnection(conn net.Conn) {
	cm.connectionsM.Lock()
	defer cm.connectionsM.Unlock()
	cm.connections[conn] = struct{}{}
}
