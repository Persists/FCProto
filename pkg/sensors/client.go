package sensors

import (
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/utils"
)

type SensorClient struct {
	sensors []struct {
		Sensor   BaseSensor
		Interval time.Duration
	}
}

func NewClient() *SensorClient {
	return &SensorClient{}
}

func (sc *SensorClient) Init() {
	sc.sensors = []struct {
		Sensor   BaseSensor
		Interval time.Duration
	}{
		{Sensor: NewCpuSensor(), Interval: 10 * time.Second},
		{Sensor: NewMemSensor(), Interval: 5 * time.Second},
		{Sensor: NewVirtualSensor(), Interval: 10 * time.Second},
	}
}

func (sc *SensorClient) SendToReceiver(stopChan <-chan bool, send func(models.Message)) {

	dataChan := make(chan string)

	// Start generating data for each sensor
	for _, sensor := range sc.sensors {
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

func (sc *SensorClient) Start(stopChan <-chan bool, send func(models.Message)) {
	sc.Init()
	sc.SendToReceiver(stopChan, send)
}
