package sensors

import (
	"encoding/json"
	"testing"
)

func TestVirtualSensor_GenerateData(t *testing.T) {
	vs := NewVirtualSensor()

	dataString := vs.GenerateData().ToString()
	data := &SensorData{}

	err := json.Unmarshal([]byte(dataString), data)

	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}

	if data.Temperature < MinTemperature || data.Temperature > MaxTemperature {
		t.Errorf("Expected temperature between %f and %f but got %f", MinTemperature, MaxTemperature, data.Temperature)
	}

	if data.Humidity < MinHumidity || data.Humidity > MaxHumidity {
		t.Errorf("Expected humidity between %f and %f but got %f", MinHumidity, MaxHumidity, data.Humidity)
	}
}
