package main

import (
	"log"

	client_config "github.com/Persists/fcproto/internal/client/client-config"
	"github.com/Persists/fcproto/internal/shared/sender"
	"github.com/Persists/fcproto/pkg/sensors"
)

type FogClient struct {
	serverConfig  *client_config.ClientConfig
	senderManager *sender.SenderManager
	sensorManager *sensors.SensorManager
}

func NewFogClient() *FogClient {
	return &FogClient{
		serverConfig:  &client_config.ClientConfig{},
		senderManager: sender.NewSenderManager(),
		sensorManager: sensors.NewSensorManager(),
	}
}

func (fc *FogClient) Init() error {
	config, err := client_config.LoadConfig()
	if err != nil {
		log.Printf("failed to load env config: %v", err)
		return err
	}
	fc.serverConfig = config

	err = fc.senderManager.Init(fc.serverConfig)
	if err != nil {
		log.Printf("failed to initialize the sender manager: %v", err)
		return err
	}

	fc.sensorManager.Init()

	return nil
}

func (fc *FogClient) Start() error {
	fc.senderManager.Start()
	fc.sensorManager.Start(fc.senderManager.StopChan, fc.senderManager.Sender.Send)

	return nil
}

func (fc *FogClient) Stop() error {
	err := fc.senderManager.Stop()
	if err != nil {
		log.Printf("failed to stop the sender manager: %v", err)
		return err
	}

	return nil
}
