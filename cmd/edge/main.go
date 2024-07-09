package main

import (
	edge "github.com/Persists/fcproto/internal/edge"
	"log"
	"time"
)

func main() {
	fc := edge.NewClient()
	err := fc.Init()
	if err != nil {
		log.Fatalf("failed to initialize the edge client: %v", err)
	}

	time.Sleep(1000 * time.Second)
}
