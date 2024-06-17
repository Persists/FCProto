package receiver

import (
	"net"

	"github.com/Persists/fcproto/internal/shared/queue"
)

type Sender struct {
	Addr  string
	queue *queue.Queue

	conn *net.Conn
}

func NewSender(addr string) *Sender {
	return &Sender{Addr: addr}
}

func (s *Sender) Connect() error {
	conn, err := net.Dial("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.queue = queue.NewQueue()

	s.conn = &conn

	return nil
}

// takes the messages from the queue and sends them to the receiver
func (s *Sender) routine() {
	for {
		msg, ok := s.queue.Dequeue()
		if !ok {
			continue
		}

		// send the message
	}
}
