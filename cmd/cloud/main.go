package main

import (
	"log"

	"github.com/Persists/fcproto/internal/cloud"
	"github.com/Persists/fcproto/internal/cloud/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("failed to load env config: %v", err)
		return
	}

	// create a new cloud client
	cc := cloud.NewClient()

	// initialize the cloud client
	err = cc.Init(config)
	if err != nil {
		log.Fatalf("failed to initialize the cloud client: %v", err)
	}

	// start the cloud client
	err = cc.Start()
	if err != nil {
		log.Fatalf("failed to start the cloud client: %v", err)
	}
}
