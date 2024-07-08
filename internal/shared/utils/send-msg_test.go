package utils

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/Persists/fcproto/internal/shared/models"
)

func TestSendMessage(t *testing.T) {
	server, client := net.Pipe()

	msg := models.Message{
		Topic: models.Sensor,
		Payload: &map[string]interface{}{
			"timestamp": 1,
			"content":   "msg1",
		},
	}

	go func() {
		err := SendMessage(client, msg)
		if err != nil {
			t.Errorf("Expected nil but got %v", err)
		}
	}()

	buf := make([]byte, 1024)
	n, err := server.Read(buf)
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}

	if n == 0 {
		t.Errorf("Expected non-zero but got %d", n)
	}

	message := &models.Message{}

	err = json.Unmarshal(buf[:n], message)
	if err != nil {
		t.Errorf("Expected nil but got %v", message)
	}

	if message.Topic != models.Sensor {
		t.Errorf("Expected sensor but got %v", message.Topic)
	}
}
