package main

import (
	"log"
	"time"

	"github.com/Persists/fcproto/internal/fog"
)

func main() {
	fc := fog.NewClient()
	err := fc.Init()
	if err != nil {
		log.Fatalf("failed to initialize the fog client: %v", err)
	}

	time.Sleep(1000 * time.Second)
}
