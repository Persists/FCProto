package main

import (
	"log"
	"time"

	"github.com/Persists/fcproto/internal/cloud"
	"github.com/Persists/fcproto/internal/cloud/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("failed to load env config: %v", err)
		return
	}

	cc := cloud.NewClient()
	cc.Init(config)

	cc.Start()

	time.Sleep(1000 * time.Second)
}
