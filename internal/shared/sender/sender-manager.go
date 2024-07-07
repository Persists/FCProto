package sender

import (
	"log"
	"time"

	client_config "github.com/Persists/fcproto/internal/client/client-config"
)

type SenderManager struct {
	Sender   *Sender
	DataChan chan string
	StopChan chan bool
}

func NewSenderManager() *SenderManager {
	return &SenderManager{
		DataChan: make(chan string),
		StopChan: make(chan bool),
	}
}

func (sm *SenderManager) Init(config *client_config.ClientConfig) error {
	sm.Sender = NewSender(config)

	_, err := sm.Sender.Connect()
	if err != nil {
		log.Printf("failed to connect: %v", err)
		return err
	}

	return nil
}

func (sm *SenderManager) Start() {
	sm.Sender.Start()
}

func (sm *SenderManager) Stop() error {
	close(sm.StopChan)

	// Give some time for goroutines to finish
	time.Sleep(1 * time.Second)
	close(sm.DataChan)

	if err := sm.Sender.Close(); err != nil {
		log.Printf("failed to close: %v", err)
		return err
	}
	return nil
}
