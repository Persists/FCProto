package connection

import (
	"fmt"
	"github.com/Persists/fcproto/internal/shared/utils"
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
				logMsg := fmt.Sprintf("Client already exists: %s", conn.RemoteAddr().String())
				log.Println(utils.Colorize(utils.Blue, logMsg))

				if listenerClient.Connections[ip].stopped() {
					listenerClient.Connections[ip].conn = &conn

					listenerClient.Connections[ip].stop = make(chan struct{})
					go listenerClient.Connections[ip].sendRoutine()
					go listenerClient.Connections[ip].receiveRoutine()
				}

				logMsg = fmt.Sprintf("Queues reattached to: %s", conn.RemoteAddr().String())
				log.Println(utils.Colorize(utils.Blue, logMsg))

				continue
			} else {
				logMsg := fmt.Sprintf("New client connected %s", conn.RemoteAddr().String())
				log.Println(utils.Colorize(utils.Blue, logMsg))
				listenerClient.Connections[ip] = newListenerConnection(&conn)

				listenerClient.Connections[ip].stop = make(chan struct{})
				go listenerClient.Connections[ip].sendRoutine()
				go listenerClient.Connections[ip].receiveRoutine()

				logMsg = fmt.Sprintf("Queue setup completed for: %s", conn.RemoteAddr().String())
				log.Println(utils.Colorize(utils.Blue, logMsg))

				go func() {
					totalMessagesFromClient := 0
					for {
						msg := listenerClient.Connections[ip].Receive()
						totalMessagesFromClient += 1
						if totalMessagesFromClient%10 == 0 {
							totalMsg := fmt.Sprintf("Total messages from %s: %d", ip, totalMessagesFromClient)
							log.Println(utils.Colorize(utils.Purple, totalMsg))
						}
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
