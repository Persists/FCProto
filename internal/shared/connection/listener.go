package connection

import (
	"log"
	"net"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

type ListenerClient struct {
	Connections map[string]*ConnectionClient
}

func newListenerClient() *ListenerClient {
	return &ListenerClient{
		Connections: make(map[string]*ConnectionClient),
	}
}

// Listen listens for incoming connections on the given socket address
// and starts goroutines to handle the incoming connections and messages
func Listen(socketAddr string, onReceive func(msg *models.Message, cc *ConnectionClient)) *ListenerClient {
	listenerClient := newListenerClient()

	listener, err := net.Listen("tcp", socketAddr)
	if err != nil {
		log.Fatalf("failed to start listener: %v", err)
	}

	go func() {
		log.Println("Listening for connections")
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %v", err)
				continue
			}

			ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
			if err != nil {
				log.Printf("failed to split host and port: %v", err)
				continue
			}

			// if connection already exists, do not create new queues
			// we just need to restart the goroutines
			if _, ok := listenerClient.Connections[ip]; ok {
				log.Printf("Connection already exists for %s", conn.RemoteAddr().String())

				if listenerClient.Connections[ip].stopped() {
					listenerClient.Connections[ip].conn = &conn

					listenerClient.Connections[ip].stop = make(chan struct{})
					go listenerClient.Connections[ip].sendRoutine()
					go listenerClient.Connections[ip].receiveRoutine()
				}

				continue
			} else {
				listenerClient.Connections[ip] = newListenerConnection(&conn)

				listenerClient.Connections[ip].stop = make(chan struct{})
				go listenerClient.Connections[ip].sendRoutine()
				go listenerClient.Connections[ip].receiveRoutine()

				go func() {
					for {
						msg := listenerClient.Connections[ip].Receive()
						// Debug print
						println("Queue Len:", listenerClient.Connections[ip].receiveQueue.Len())
						onReceive(&msg, listenerClient.Connections[ip])
					}
				}()

			}

		}
	}()

	return listenerClient
}

func newListenerConnection(conn *net.Conn) *ConnectionClient {
	return &ConnectionClient{
		receiveQueue: queue.New[models.Message](),
		sendQueue:    queue.New[models.Message](),

		conn: conn,
	}

}

func (lc *ListenerClient) BroadcastMsg(msg models.Message) {
	for _, cc := range lc.Connections {
		cc.Send(msg)
	}
}

func (lc *ListenerClient) GetAllIPs() []string {
	var ips []string
	for ip := range lc.Connections {
		ips = append(ips, ip)
	}
	return ips
}
