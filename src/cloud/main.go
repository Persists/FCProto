package main

import (
	"time"

	"github.com/Persists/fcproto/internal/shared/connection"
	"github.com/Persists/fcproto/internal/shared/models"
)

func onReceive(message *models.Message, cc *connection.ConnectionClient) {
	println(message.Topic, message.Payload)
}

func main() {

	_ = connection.Listen("localhost:8080", onReceive)

	time := time.NewTimer(1000 * time.Second)

	<-time.C
}
