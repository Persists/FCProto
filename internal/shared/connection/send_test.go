package connection

import (
	"net"
	"testing"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

func TestConnectionClient_Send(t *testing.T) {

	server, client := net.Pipe()

	ccClient := &ConnectionClient{
		ingress: &queue.Queue[models.Message]{},
		egress:    &queue.Queue[models.Message]{},

		stop: make(chan struct{}),

		conn: &client,
	}

	ccServer := &ConnectionClient{
		ingress: &queue.Queue[models.Message]{},
		egress:  &queue.Queue[models.Message]{},

		stop: make(chan struct{}),

		conn: &server,
	}

	go ccClient.sendRoutine()
	go ccClient.receiveRoutine()
	go ccServer.sendRoutine()
	go ccServer.receiveRoutine()

	message := models.Message{
		Topic:   "test",
		Payload: &map[string]interface{}{"message": "Hello, World!"},
	}

	ccClient.Send(message)

	if ccClient.egress.Len() != 0 {
		t.Errorf("expected sendQueue should not be empty, got %d", ccClient.egress.Len())
	}

	msg := ccServer.Receive()
	if msg.Topic != message.Topic {
		t.Errorf("expected topic to be %s, got %s", message.Topic, msg.Topic)
	}

}
