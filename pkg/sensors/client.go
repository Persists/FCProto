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
		{Sensor: NewCpuSensor(), Interval: 4 * time.Second},
		{Sensor: NewMemSensor(), Interval: 4 * time.Second},
		{Sensor: NewVirtualSensor(), Interval: 2 * time.Second},
	}
}

// SendToReceiver starts ticker for each sensor to generate sensor data to be sent to the cloud application receiver
func (sc *SensorClient) SendToReceiver(stopChan <-chan bool, send func(models.Message)) {

	dataChan := make(chan string)

	for _, sensor := range sc.sensors {
		go utils.StartTicker(sensor.Interval, func() string {
			return sensor.Sensor.GenerateData().ToString()
		}, stopChan, dataChan)
	}

	go func() {
		for data := range dataChan {
			send(models.NewMessage(models.Sensor, models.NewSensorMessage(data)))
		}
	}()
}

// Start initializes the sensors and starts the data generation process
func (sc *SensorClient) Start(stopChan <-chan bool, send func(models.Message)) {
	sc.Init()
	sc.SendToReceiver(stopChan, send)
}
