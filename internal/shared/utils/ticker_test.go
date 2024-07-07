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

func TestMultipleTickers(t *testing.T) {
	// Test parameters
	interval := 100 * time.Millisecond
	testDuration := 1050 * time.Millisecond
	expectedTicks := int(testDuration / interval)
	counter1 := 0
	counter2 := 0

	// Mock tick function
	mockTickFunc := func() int {
		return 1
	}

	stopChan := make(chan bool)
	// Channels for communication
	dataTickerChan1 := make(chan int)
	dataTickerChan2 := make(chan int)

	// Start the tickers in separate goroutines
	go StartTicker(interval, mockTickFunc, stopChan, dataTickerChan1)
	go StartTicker(interval, mockTickFunc, stopChan, dataTickerChan2)

	go func() {
		for range dataTickerChan1 {
			counter1++
		}
	}()

	go func() {
		for range dataTickerChan2 {
			counter2++
		}
	}()

	time.Sleep(testDuration)

	close(stopChan)

	time.Sleep(1 * time.Second)

	if counter1 != expectedTicks {
		t.Errorf("Expected %d ticks for ticker 1 but got %d", expectedTicks, counter1)
	}

	if counter2 != expectedTicks {
		t.Errorf("Expected %d ticks for ticker 2 but got %d", expectedTicks, counter2)
	}
}
