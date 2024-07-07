package sensors

import (
	"testing"
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
)

func TestGenerateDataAtInterval(t *testing.T) {

	mockSensor := &MockSensor{}
	mockManager := NewSensorManager()

	mockManager.sensors = append(mockManager.sensors, mockSensor)

	sendMessages := []models.Message{}

	send := func(msg models.Message) {
		sendMessages = append(sendMessages, msg)
	}

	stopChan := make(chan bool)

	mockManager.SendToReceiver(stopChan, send)
	mockManager.SendToReceiver(stopChan, send)

	time.Sleep(5 * time.Second)

	close(stopChan)

	time.Sleep(2 * time.Second)

	if len(sendMessages) == 0 {
		t.Errorf("Expected to receive messages but got none")
	}

	if len(sendMessages) != 4 {
		t.Errorf("Expected to receive 2 messages but got %d", len(sendMessages))
	}
}

type MockSensor struct{}

func (m *MockSensor) GenerateData() string {
	return "mock data"
}
