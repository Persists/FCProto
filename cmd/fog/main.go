package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/Persists/fcproto/internal/shared/models"
)

const (
	addr = "localhost:5555"
)

func main() {
	sender(addr)
}

func sender(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	log.Printf("connected to %s", address)

	msg := models.Message{
		Timestamp: 1234567890,

		Content: "Hello, World!",
	}

	for i := 0; i < 10; i++ {
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("failed to marshal message: %v", err)
		}

		n, err := conn.Write([]byte(jsonMsg))
		if err != nil {
			log.Fatalf("failed to write to connection: %v", err)
		}

		log.Printf("sent %d bytes: %s", n, jsonMsg)
	}
}
