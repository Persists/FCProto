package connection

import (
	"log"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
)

// this routine is used to pop messages from the send queue
// and send them to the connection subsequently
func (cc *ConnectionClient) sendRoutine() {
	for {
		message := cc.egress.Dequeue()

		if stop := cc.stopped(); stop {
			cc.egress.Enqueue(message)
			return
		}

		err := utils.SendMessage(*cc.conn, message)

		if err != nil {
			log.Println(utils.Colorize(utils.Red, "Failed to send message"))
		}
	}
}

// Send enqueues a message to the send queue
func (cc *ConnectionClient) Send(message models.Message) {
	cc.egress.Enqueue(message)
}
