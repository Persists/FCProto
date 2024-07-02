package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Persists/fcproto/internal/shared/models"
	"log"
	"net"
)

const defaultPort = 54548

func sendMessage(conn net.Conn, msg models.Message) ([]byte, error) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
	}

	n, err := conn.Write([]byte(jsonMsg))
	if err != nil {
		log.Printf("failed to write to conn: %v", err)
	}

	log.Printf("sent %d bytes", n)
	return jsonMsg, err
}

// CheckPortAvailable checks if a port is available for use
func checkPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func EstablishConnection(socketAddr string) (net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)
	if checkPortAvailable(defaultPort) {
		dialer := net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP:   net.IPv4zero,
				Port: defaultPort,
			},
		}

		conn, err = dialer.Dial("tcp", socketAddr)
	} else {
		fmt.Printf("Default port %d is not available, using random one\n", defaultPort)
		conn, err = net.Dial("tcp", socketAddr)
	}

	if err != nil {
		fmt.Print("Failed to connect to server for heartbeat: ", err)
	}
	return conn, err
}
