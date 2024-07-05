package utils

import (
	"testing"
	"time"
)

func TestStartTicker(t *testing.T) {
	// Test parameters
	interval := 100 * time.Millisecond
	testDuration := 1050 * time.Millisecond
	expectedTicks := int(testDuration / interval)
	counter := 0

	// Mock tick function
	mockTickFunc := func() int {
		return 1
	}

	// Channels for communication
	stopTickerChan := make(chan bool)
	dataTickerChan := make(chan int)

	// Start the ticker in a separate goroutine
	go StartTicker(interval, mockTickFunc, stopTickerChan, dataTickerChan)

	go func() {
		for range dataTickerChan {
			counter++
		}
	}()

	time.Sleep(testDuration)

	stopTickerChan <- true

	if counter != expectedTicks {
		t.Errorf("Expected %d ticks but got %d", expectedTicks, counter)
	}

}
