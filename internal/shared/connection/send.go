package connection

import (
	"fmt"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
)

// this routine is used to pop messages from the send queue
// and send them to the connection subsequently
func (cc *ConnectionClient) sendRoutine() {
	for {
		message := cc.sendQueue.Dequeue()

		if stop := cc.stopped(); stop {
			cc.sendQueue.Enqueue(message)
			return
		}

		err := utils.SendMessage(*cc.conn, message)

		if err != nil {
			fmt.Println("Failed to send message")
		}
	}
}

// Send enqueues a message to the send queue
func (cc *ConnectionClient) Send(message models.Message) {
	cc.sendQueue.Enqueue(message)
}
