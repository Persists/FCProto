package sender

import (
	"testing"

	client_config "github.com/Persists/fcproto/internal/client/client-config"
	"github.com/Persists/fcproto/internal/shared/models"
)

func TestSenderManager_Init(t *testing.T) {
	config := &client_config.ClientConfig{
		BaseEnv: &models.BaseEnv{SocketAddr: "localhost:8080"},
	}

	senderManager := NewSenderManager()
	err := senderManager.Init(config)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if senderManager.Sender == nil {
		t.Errorf("Expected sender to be initialized, but it was nil")
	}
}

func TestSenderManager_Start(t *testing.T) {
	config := &client_config.ClientConfig{
		BaseEnv: &models.BaseEnv{SocketAddr: "localhost:8080"},
	}

	senderManager := NewSenderManager()
	err := senderManager.Init(config)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	senderManager.Start()

}

func TestSenderManager_Stop(t *testing.T) {
	config := &client_config.ClientConfig{
		BaseEnv: &models.BaseEnv{SocketAddr: "localhost:8080"},
	}

	senderManager := NewSenderManager()
	err := senderManager.Init(config)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	senderManager.Start()
	err = senderManager.Stop()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Check if the channels are closed
	select {
	case _, ok := <-senderManager.DataChan:
		if ok {
			t.Errorf("Expected DataChan to be closed, but it was open")
		}
	default:
	}

	select {
	case _, ok := <-senderManager.StopChan:
		if ok {
			t.Errorf("Expected StopChan to be closed, but it was open")
		}
	default:
	}
}

func TestSender_Send(t *testing.T) {
	config := &client_config.ClientConfig{
		BaseEnv: &models.BaseEnv{SocketAddr: "localhost:8080"},
	}

	sender := NewSender(config)
	_, err := sender.Connect()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	message := models.NewMessage(models.Sensor, models.NewSensorMessage("test data"))
	sender.Send(message)

}
