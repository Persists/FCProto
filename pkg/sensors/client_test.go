package sensors

import (
	"testing"
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
)

func TestGenerateDataAtInterval(t *testing.T) {

	mockSensor := &MockSensor{}
	mockClient := NewClient()

	sensors := []struct {
		Sensor   BaseSensor
		Interval time.Duration
	}{
		{Sensor: mockSensor, Interval: 1 * time.Second},
		{Sensor: mockSensor, Interval: 2 * time.Second},
	}

	mockClient.sensors = append(mockClient.sensors, sensors...)

	sendMessages := []models.Message{}

	send := func(msg models.Message) {
		sendMessages = append(sendMessages, msg)
	}

	stopChan := make(chan bool)

	mockClient.SendToReceiver(stopChan, send)

	time.Sleep(4 * time.Second)

	close(stopChan)

	time.Sleep(2 * time.Second)

	if len(sendMessages) == 0 {
		t.Errorf("Expected to receive messages but got none")
	}

	if len(sendMessages) != 6 {
		t.Errorf("Expected to receive 4 messages but got %d", len(sendMessages))
	}
}

type MockSensor struct{}

func (m *MockSensor) GenerateData() *BaseSensorData {
	return &BaseSensorData{
		Data: "Mock Sensor Data",
	}

}
