package receiver

import (
	"encoding/json"
	"log"
	"net"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

type Sender struct {
	Addr  string
	queue *queue.Queue[models.Message]

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

	s.queue = &queue.Queue[models.Message]{}

	s.conn = &conn

	s.routine()

	return nil
}

func (s *Sender) Send(msg models.Message) {
	s.queue.Enqueue(msg)
}

// takes the messages from the queue and sends them to the receiver
func (s *Sender) routine() {
	go func() {
		for {
			msg := s.queue.Dequeue()

			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			if _, err := (*s.conn).Write(data); err != nil {
				log.Printf("Error sending message: %v", err)
			}

		}
	}()

}
