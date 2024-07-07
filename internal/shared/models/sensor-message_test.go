package models

import (
	"testing"
	"time"
)

func TestNewSensorMessage(t *testing.T) {
	content := "msg1"
	sensorMessage := NewSensorMessage(content)

	if sensorMessage.Content != content {
		t.Errorf("Expected %s, got %s", content, sensorMessage.Content)
	}

	time.Sleep(100 * time.Millisecond)

	if sensorMessage.Timestamp > time.Now().Unix() {
		t.Errorf("Expected timestamp to be in the past, got %d", sensorMessage.Timestamp)
	}
}
