package models

import "testing"

func TestNewMessage(t *testing.T) {
	topic := Sensor
	payload := SensorMessage{
		Timestamp: 1,
		Content:   "msg1",
	}

	message := NewMessage(topic, payload)

	if message.Topic != topic {
		t.Errorf("Expected %v but got %v", topic, message.Topic)
	}

	if message.Payload == nil {
		t.Errorf("Expected non-nil but got nil")
	}
}
