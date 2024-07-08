package connection

import (
	"net"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

type ConnectionClient struct {
	stop chan struct{}

	receiveQueue *queue.Queue[models.Message]
	sendQueue    *queue.Queue[models.Message]

	conn *net.Conn
}

func (cc *ConnectionClient) RemoteAddress() string {
	conn := *cc.conn

	return conn.RemoteAddr().String()
}

func newConnectionClient() *ConnectionClient {
	return &ConnectionClient{
		receiveQueue: queue.New[models.Message](),
		sendQueue:    queue.New[models.Message](),
	}
}

func (cc *ConnectionClient) stopped() bool {
	select {
	case <-cc.stop:
		return true
	default:
		return false
	}
}
