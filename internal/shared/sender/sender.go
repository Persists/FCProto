package sender

import (
	"encoding/json"
	"fmt"
	client_config "github.com/Persists/fcproto/internal/client/client-config"
	sharedUtils "github.com/Persists/fcproto/internal/shared/utils"
	"log"
	"net"
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

type Sender struct {
	queue *queue.Queue[models.Message]
	conn  *net.Conn

	*client_config.ClientConfig
}

func NewSender(config *client_config.ClientConfig) *Sender {
	return &Sender{
		queue: queue.NewQueue[models.Message](),

		ClientConfig: config,
	}
}

func (s *Sender) Connect() (*net.Conn, error) {
	retryDelay := 2 * time.Second
	maxRetryDelay := 60 * time.Second // Maximum retry delay
	retries, maxRetries := 0, 1       // Maximum number of retries

	for {
		if retries < maxRetries {
			conn, err := sharedUtils.EstablishConnection(s.SocketAddr, s.SendPort)
			if err != nil {
				retries++
				log.Printf("Failed to connect to %s: %v. Retrying in %s...", s.SocketAddr, err, retryDelay)
				time.Sleep(retryDelay)
				retryDelay *= 2
				if retryDelay > maxRetryDelay {
					retryDelay = maxRetryDelay
				}
				continue
			} else {
				log.Printf("Connected to %s", s.SocketAddr)
				s.conn = &conn
				err := s.sendInitialHeartbeat()
				if err != nil {
					fmt.Printf("Failed to send initial heartbeat: %v", err)
					return nil, err
				}
				return s.conn, nil
			}
		} else {
			log.Printf("Failed to connect to %s after maximum retries. Hosting callback Port", s.SocketAddr)
			_, err := s.startCallbackListener()
			if err != nil {
				fmt.Printf("Failed to start callback listener: %v", err)
				return nil, err
			}
			retries = 0
			continue
		}
	}
}

func (s *Sender) Start() {
	s.routine()
}

func (s *Sender) Close() error {
	if err := (*s.conn).Close(); err != nil {
		return err
	}
	return nil
}

func (s *Sender) Send(msg models.Message) {
	s.queue.Enqueue(msg)
}

// takes the messages-utils from the queue and sends them to the sender
func (s *Sender) routine() {
	go func() {
		for {
			msg := s.queue.Dequeue()

			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			err = s.writeWithRetry(data)
			if err != nil {
				fmt.Printf("Failed to send message: %v", err)
				return
			}
		}
	}()
}

func (s *Sender) writeWithRetry(data []byte) error {
	for {
		_, err := (*s.conn).Write(data)
		if err != nil {
			log.Printf("Error sending data: %v", err)
			_, err := s.Connect()
			if err != nil {
				fmt.Printf("Failed to reconnect: %v", err)
				return err
			}
			return err
		}
		return nil
	}
}

func (s *Sender) sendInitialHeartbeat() error {
	msg := models.NewMessage(models.Heartbeat, models.NewHeartbeatMessage(s.NotifyAddr))
	fmt.Println(s.NotifyAddr)
	_, err := sharedUtils.SendMessage(*s.conn, msg)
	if err != nil {
		fmt.Printf("Failed to send initial heartbeat: %v", err)
		return err
	}
	return nil
}

func (s *Sender) startCallbackListener() (conn net.Conn, err error) {
	listener, err := net.Listen("tcp", "fog:5556")
	if err != nil {
		log.Fatalf("Failed to start listener on port %s: %v", s.NotifyAddr, err)
	}
	defer func(listener net.Listener) {
		closeErr := listener.Close()
		if err != nil {
			fmt.Printf("Failed to close listener: %v", err)
			err = closeErr
		}
	}(listener)
	log.Printf("Listening on default port %s", s.NotifyAddr)

	conn, err = listener.Accept()
	if err != nil {
		log.Printf("Failed to accept connection: %v", err)
		return
	}
	return
}
