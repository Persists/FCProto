package sensors

import (
	"encoding/json"

	"github.com/jaswdr/faker"
)

const (
	MaxDecimals    = 4
	MinHumidity    = 30.0
	MaxHumidity    = 50.0
	MinTemperature = -10.0
	MaxTemperature = 30.0
)

// SensorData holds the simulated data from the sensors.
type SensorData struct {
	Temperature float64
	Humidity    float64
}

// VirtualSensor simulates a sensor.
type VirtualSensor struct {
	faker faker.Faker

	BaseSensor
}

// NewVirtualSensor creates a new VirtualSensor.
func NewVirtualSensor() *VirtualSensor {
	return &VirtualSensor{
		faker: faker.New(),
	}
}

// GenerateFakeData generates simulated sensor data.
func (vs VirtualSensor) GenerateFakeData() SensorData {
	temperature := vs.faker.Float64(MaxDecimals, MinTemperature, MaxTemperature)
	humidity := vs.faker.Float64(MaxDecimals, MinHumidity, MaxHumidity)
	return SensorData{
		Temperature: temperature,
		Humidity:    humidity,
	}
}

func (vs VirtualSensor) GenerateData() (dataString string) {
	fakeData := vs.GenerateFakeData()
	dataBytes, err := json.Marshal(fakeData)
	if err != nil {
		dataString = ""
	} else {
		dataString = string(dataBytes)
	}
	return
}
