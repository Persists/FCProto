package sensors

import (
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
)

type SensorManager struct {
	sensors []struct {
		Sensor   BaseSensor
		Interval time.Duration
	}
}

func NewSensorManager() *SensorManager {
	return &SensorManager{}
}

func (manager *SensorManager) Init() {
	manager.sensors = []struct {
		Sensor   BaseSensor
		Interval time.Duration
	}{
		{Sensor: NewCpuSensor(), Interval: 1 * time.Second},
		{Sensor: NewMemSensor(), Interval: 2 * time.Second},
		{Sensor: NewVirtualSensor(), Interval: 3 * time.Second},
	}
}

func (manager *SensorManager) SendToReceiver(stopChan <-chan bool, send func(models.Message)) {

	dataChan := make(chan string)

	// Start generating data for each sensor
	for _, sensor := range manager.sensors {
		go utils.StartTicker(sensor.Interval, func() string {
			return sensor.Sensor.GenerateData().ToString()
		}, stopChan, dataChan)
	}

	// Collect and print data for demonstration purposes
	go func() {
		for data := range dataChan {
			send(models.NewMessage(models.Sensor, models.NewSensorMessage(data)))
		}
	}()
}

func (manager *SensorManager) Start(stopChan <-chan bool, send func(models.Message)) {
	manager.Init()
	manager.SendToReceiver(stopChan, send)
}
