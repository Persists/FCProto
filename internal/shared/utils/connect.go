package utils

import (
	"fmt"
	"net"
)

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

// EstablishConnection establishes a connection to a server using the given socket address and port
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
