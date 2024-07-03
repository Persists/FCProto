package main

import (
	"log"
	"time"
)

func main() {
	fogClient := NewFogClient()
	err := fogClient.Init()
	if err != nil {
		log.Fatalf("failed to initialize the fog client: %v", err)
	}

	err = fogClient.Start()
	if err != nil {
		log.Fatalf("failed to start the fog client: %v", err)
	}

	time.Sleep(30 * time.Second)

	err = fogClient.Stop()
	if err != nil {
		log.Fatalf("failed to stop the fog client: %v", err)
	}
}
