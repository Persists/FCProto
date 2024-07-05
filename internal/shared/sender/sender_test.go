package sender

import (
	"net"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

type SpySender struct {
	connectCalled bool
	closeCalled   bool

	messages []models.Message
	conn     *net.Conn
	queue    *queue.Queue[models.Message]
}

func test() {
	s := Sender{SpySender{}}
}

func (s *SpySender) Connect() (*net.Conn, error) {
	s.connectCalled = true
	return s.conn, nil
}

func (s *SpySender) Close() error {
	s.closeCalled = true
	if s.conn != nil {
		return (*s.conn).Close()
	}
	return nil
}

func (s *SpySender) Send(msg models.Message) {
	s.messages = append(s.messages, msg)
	s.queue.Enqueue(msg)
}

func (s *SpySender) Start() {
	s.routine()
}

func (s *SpySender) routine() {
}

func (s *SpySender) writeWithRetry() error {
	return nil
}

func (s *SpySender) sendInitialHeartbeat() error {
	return nil
}

func (s *SpySender) startCallbackListener() (conn net.Conn, err error) {
	return nil, nil
}
