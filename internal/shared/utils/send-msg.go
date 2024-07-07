package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Persists/fcproto/internal/shared/models"
	"log"
	"net"
)

func SendMessage(conn net.Conn, msg models.Message) ([]byte, error) {
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
	err = ln.Close()
	if err != nil {
		fmt.Printf("Failed to close listener: %v", err)
		return false
	}
	return true
}

func EstablishConnection(socketAddr string, sendPort int) (net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)
	if checkPortAvailable(sendPort) {
		dialer := net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP:   net.IPv4zero,
				Port: sendPort,
			},
		}

		conn, err = dialer.Dial("tcp", socketAddr)
	} else {
		fmt.Printf("Default port %d is not available, using random one\n", sendPort)
		conn, err = net.Dial("tcp", socketAddr)
	}

	if err != nil {
		fmt.Print("Failed to connect to server: ", err)
	}
	return conn, err
}
