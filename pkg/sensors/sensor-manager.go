package sensors

import (
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/sender"
	"github.com/Persists/fcproto/internal/shared/utils"
	"time"
)

type SensorManager struct {
	sensors []BaseSensor
}

func NewSensorManager() *SensorManager {
	return &SensorManager{}
}

func (manager *SensorManager) Init() {
	// Hardware Sensor
	/*
		hwSensor, err := sensors.GetSystemInfo()
		if err != nil {
			log.Fatalf("Could not get system information for hwSensor %v\n", err)
		}
		sensorsArray = append(sensorsArray, hwSensor)
	*/

	// Virtual Sensor
	manager.sensors = append(manager.sensors, NewVirtualSensor())
}

func (manager *SensorManager) SendToReceiver(sm *sender.SenderManager) {
	// Start generating data for each sensor
	for _, sensor := range manager.sensors {
		go utils.StartTicker(2*time.Second, func() string {
			return sensor.GenerateData()
		}, sm.StopChan, sm.DataChan)
	}

	// Collect and print data for demonstration purposes
	go func() {
		for data := range sm.DataChan {
			sm.Sender.Send(models.NewMessage(models.Sensor, models.NewSensorMessage(data)))
		}
	}()
}
