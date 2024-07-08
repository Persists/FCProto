package utils

import (
	"net"
	"strings"
	"testing"
)

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

	splitted := strings.Split(conn.LocalAddr().String(), ":")
	port := splitted[len(splitted)-1]

	if port != "8081" {
		t.Errorf("Expected 8081 but got %s", port)
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

	splitted := strings.Split(conn.LocalAddr().String(), ":")
	port := splitted[len(splitted)-1]

	if port == "8080" {
		t.Errorf("Expected random port but got %s", port)
	}
}
