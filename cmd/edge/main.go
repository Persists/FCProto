package main

import (
	"log"

	edge "github.com/Persists/fcproto/internal/edge"
)

func main() {
	// create a new edge client
	fc := edge.NewClient()

	// initialize the edge client
	err := fc.Start()
	if err != nil {
		log.Fatalf("failed to initialize the edge client: %v", err)
	}

}
