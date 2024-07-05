package utils

import (
	"log"
	"time"
)

// StartTicker starts a ticker that executes a generic function at each tick interval
func StartTicker[T any](interval time.Duration, tickFunc func() T, stopTickerChan <-chan bool, dataTickerChan chan<- T) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dataTickerChan <- tickFunc()
		case <-stopTickerChan:
			log.Println("Stopping ticker")
			return
		}
	}
}
