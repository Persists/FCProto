package connection

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"github.com/Persists/fcproto/internal/shared/models"
)

// this routine is used to receive messages from the connection and
// enqueue them to the receive queue ( this routine is responsible for
// tracking the liveness of the connection)
func (cc *ConnectionClient) receiveRoutine() {

	reader := bufio.NewReader(*cc.conn)
	decoder := json.NewDecoder(reader)

	for {
		var message models.Message
		err := decoder.Decode(&message)

		// if the connection is closed, we stop the routine
		// and emit a signal to other routines to stop
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed to %s", (*cc.conn).RemoteAddr())

				close(cc.stop)
				return
			}
			log.Printf("Failed to decode JSON: %v", err)
			continue
		}

		cc.ingress.Enqueue(message)
	}
}

// Receive returns the next message in the receive queue
// if no messages are in the queue it blocks until a message is received
func (cc *ConnectionClient) Receive() models.Message {
	return cc.ingress.Dequeue()
}
