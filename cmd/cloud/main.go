package main

import (
	"io"
	"log"
	"net"
)

const (
	addr = "localhost:5555"
)

func main() {
	server(addr)

}

func server(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("listening on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("failed to accept: %v", err)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		const bufSize = 1024
		buf := make([]byte, bufSize)

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client: %v", conn.RemoteAddr())
				return
			}
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		log.Printf("Received %d bytes: %s", n, string(buf[:n]))
	}
}
