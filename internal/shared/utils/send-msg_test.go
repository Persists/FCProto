package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
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
		_, err := SendMessage(client, msg)
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

func TestCheckPortAvailable(t *testing.T) {
	if !checkPortAvailable(8080) {
		t.Errorf("Expected true but got false")
	}

	conn, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}

	defer conn.Close()

	if checkPortAvailable(8080) {
		t.Errorf("Expected false but got true")
	}

	conn.Close()

}

func TestEstablishConnection(t *testing.T) {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Errorf("Error while mocking listener: %v", err)
	}

	defer listen.Close()
	conn, err := EstablishConnection("localhost:8080", 8081)
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Errorf("Expected non-nil but got nil")
	}

	if strings.Split(conn.LocalAddr().String(), ":")[3] != "8081" {
		t.Errorf("Expected 8081 but got %v", strings.Split(conn.LocalAddr().String(), ":")[3])
	}
}

func TestEstablishConnectionRandomPort(t *testing.T) {
	_, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Errorf("Error while mocking listener: %v", err)
	}

	conn, err := EstablishConnection("localhost:8080", 8080)
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Errorf("Expected non-nil but got nil")
	}

	fmt.Println(conn.LocalAddr().String())

	if strings.Split(conn.LocalAddr().String(), ":")[3] == "8080" {
		t.Errorf("Expected random port but got 8080")
	}
}
