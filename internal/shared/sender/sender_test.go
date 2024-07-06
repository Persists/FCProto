package sender

import (
<<<<<<< HEAD
	"github.com/Persists/fcproto/internal/shared/models"
	"testing"
)

// MockQueue is a mock implementation of Queue.
type MockQueue struct {
	Messages []models.Message
}

func (q *MockQueue) Enqueue(msg models.Message) {
	q.Messages = append(q.Messages, msg)
}

func (q *MockQueue) Dequeue() models.Message {
	if len(q.Messages) == 0 {
		return models.Message{}
	}
	msg := q.Messages[0]
	q.Messages = q.Messages[1:]
	return msg
}

// MockUtils is a mock implementation of Utils.
type MockUtils struct{}

func TestSender_Send(t *testing.T) {
	mockQueue := &MockQueue{}
	mockConfig := &MockConfig{}
	mockUtils := &MockUtils{}

	sender := NewSender(mockQueue, mockConfig, mockUtils)

	msg := models.Message{
		// Initialize message as needed for test.
	}

	sender.Send(msg)

	if len(mockQueue.Messages) != 1 {
		t.Errorf("Expected queue length 1, got %d", len(mockQueue.Messages))
	}

	// Additional assertions as needed.
}

func TestSender_Connect(t *testing.T) {
	mockQueue := &MockQueue{}
	mockConfig := &MockConfig{}
	mockUtils := &MockUtils{}

	sender := NewSender(mockQueue, mockConfig, mockUtils)

	err := sender.Connect()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Additional assertions as needed.
=======
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

>>>>>>> ff03ac9 (fix: docker deployment)
}
