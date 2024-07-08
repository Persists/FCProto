package utils

import (
	"encoding/json"
	"log"
	"net"

	"github.com/Persists/fcproto/internal/shared/models"
)

func SendMessage(conn net.Conn, msg models.Message) error {
	jsonMsg, err := json.Marshal(msg)

	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return err
	}

	_, err = conn.Write(jsonMsg)

	if err != nil {
		log.Printf("failed to write message: %v", err)
		return err
	}

	return nil
}
