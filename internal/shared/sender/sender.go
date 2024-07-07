package sender

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	client_config "github.com/Persists/fcproto/internal/client/client-config"
	sharedUtils "github.com/Persists/fcproto/internal/shared/utils"

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

	// This implementats an exponential backoff algorithm
	// when it isnt able to connect to the server it will start hosting a callback port
	// (callback port: port that the server can connect to when it is ready to accept the connection)
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

func (s *Sender) Start(StopChan chan bool) {
	s.routine()
	s.printServerMessage(StopChan)
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

func (s *Sender) printServerMessage(stopChan chan bool) {
	go func() {
		reader := bufio.NewReader(*s.conn)
		decoder := json.NewDecoder(reader)

		for {
			select {
			case <-stopChan:
				return
			default:
				var msg models.Message
				err := decoder.Decode(&msg)
				if err != nil {
					if err == io.EOF {
						log.Println("Connection closed by server")
						return
					}
					log.Printf("Error decoding JSON: %v", err)
					continue
				}
				fmt.Printf("Received message from server - Topic: %s, Payload: %v\n", msg.Topic, msg.Payload)
			}
		}
	}()
}

func (s *Sender) sendInitialHeartbeat() error {
	msg := models.NewMessage(models.Heartbeat, models.NewHeartbeatMessage(s.NotifyPort))
	fmt.Println(s.NotifyPort)
	_, err := sharedUtils.SendMessage(*s.conn, msg)
	if err != nil {
		fmt.Printf("Failed to send initial heartbeat: %v", err)
		return err
	}
	return nil
}

// starts a listener on the default port and returns the connection
// if the cloud connects to the callback port
func (s *Sender) startCallbackListener() (conn net.Conn, err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.NotifyPort))
	if err != nil {
		log.Fatalf("Failed to start listener on port %s: %v", s.NotifyPort, err)
	}
	defer func(listener net.Listener) {
		closeErr := listener.Close()
		if err != nil {
			fmt.Printf("Failed to close listener: %v", err)
			err = closeErr
		}
	}(listener)
	log.Printf("Listening on default port %s", s.NotifyPort)

	conn, err = listener.Accept()
	if err != nil {
		log.Printf("Failed to accept connection: %v", err)
		return
	}
	return
}
