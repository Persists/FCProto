package connection

import (
	"log"
	"net"
	"time"
)

// Connect will connect to the given socket address and return a ConnectionClient,
// it also starts the send and receive routines, when connection is closed it will
// automatically reconnect
func Connect(socketAddress string) *ConnectionClient {
	cc := newConnectionClient()

	go func() {
	connect:
		err := cc.connectWithExpotentialBackoff(socketAddress)
		if err != nil {
			log.Printf("Failed to reconnect: %v", err)
		}

		cc.stop = make(chan struct{})
		go cc.sendRoutine()
		go cc.receiveRoutine()

		<-cc.stop
		log.Printf("Connection closed. Reconnecting...")
		goto connect
	}()

	return cc
}

// connectWithExpotentialBackoff will connect to the given socket address and implements a backoff mechanism
func (cc *ConnectionClient) connectWithExpotentialBackoff(socketAddress string) error {
	retryDelay := 2 * time.Second
	maxRetryDelay := 60 * time.Second // Maximum retry delay

	for {
		conn, err := net.Dial("tcp", socketAddress)
		if err != nil {
			log.Printf("Failed to connect to %s: %v. Retrying in %s...", socketAddress, err, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2
			if retryDelay > maxRetryDelay {
				retryDelay = maxRetryDelay
			}
			continue
		} else {
			log.Printf("Connected to %s", socketAddress)
			cc.conn = &conn
			return nil
		}

	}
}
