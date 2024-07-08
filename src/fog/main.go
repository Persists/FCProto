package main

import (
	"fmt"
	"time"

	"github.com/Persists/fcproto/internal/shared/connection"
	"github.com/Persists/fcproto/internal/shared/models"
)

func main() {
	client := connection.Connect(":8080")

	fmt.Println("Connected")

	go func() {
		for {
			msg := client.Receive()
			println(msg.Topic, msg.Payload)
		}
	}()

	time.Sleep(1 * time.Second)

	var msg = models.Message{
		Topic:   "test",
		Payload: &map[string]interface{}{"message": "Hello, World!"},
	}

	fmt.Println("Sending message")

	client.Send(msg)

	time.Sleep(5 * time.Second)

	for i := 0; i < 1000; i++ {
		msg := models.Message{
			Topic:   "test",
			Payload: &map[string]interface{}{"message": fmt.Sprintf("Hello, World! %d", i)},
		}

		time.Sleep(1 * time.Millisecond)

		client.Send(msg)
	}

	fmt.Println("Sent messages")

	time.Sleep(1000 * time.Second)
}
