package connection

import (
	"net"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

type ConnectionClient struct {
	stop chan struct{}

	ingress *queue.Queue[models.Message]
	egress  *queue.Queue[models.Message]

	conn *net.Conn
}

func (cc *ConnectionClient) RemoteAddress() string {
	conn := *cc.conn

	return conn.RemoteAddr().String()
}

func newConnectionClient() *ConnectionClient {
	return &ConnectionClient{
		ingress: queue.New[models.Message](),
		egress:  queue.New[models.Message](),
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
