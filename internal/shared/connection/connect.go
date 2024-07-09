package connection

import (
	"fmt"
	"github.com/Persists/fcproto/internal/shared/utils"
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

		stopPrint := make(chan struct{})
		utils.LogQueuePeriodically(cc.sendQueue, stopPrint)

		cc.connectWithExpotentialBackoff(socketAddress)

		close(stopPrint)

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
func (cc *ConnectionClient) connectWithExpotentialBackoff(socketAddress string) {
	retryDelay := 1 * time.Second
	maxRetryDelay := 60 * time.Second

	for {
		conn, err := net.Dial("tcp", socketAddress)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to %s: %v. Retrying in %s...", socketAddress, err, retryDelay)
			log.Println(utils.Colorize(utils.Red, msg))
			time.Sleep(retryDelay)
			retryDelay *= 2
			if retryDelay > maxRetryDelay {
				retryDelay = maxRetryDelay
			}
			continue
		} else {
			msg := fmt.Sprintf("Connected to %s", socketAddress)
			log.Println(utils.Colorize(utils.Green, msg))
			cc.conn = &conn
			return
		}

	}
}
