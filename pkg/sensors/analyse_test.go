package sensors

import (
	"testing"
)

func TestAnalyse(t *testing.T) {
	memoryData := []float64{
		9, 10, 12, 13, 13, 13, 15, 15, 16, 16, 18, 22, 23, 24, 24, 25,
	}

	avg, dev, mean, min, max := Analyse(memoryData)

	if avg != 16.75 {
		t.Errorf("Expected average to be 16.75 but got %f", avg)
	}

	if int(dev*100000) != int(5.11737*100000) {
		t.Errorf("Expected standard deviation to be 5.11737 but got %f", dev)
	}

	if mean != 16.75 {
		t.Errorf("Expected mean to be 16.75 but got %f", mean)
	}

	if min != 9 {
		t.Errorf("Expected min to be 9 but got %f", min)
	}

	if max != 25 {
		t.Errorf("Expected max to be 25 but got %f", max)
	}
}

func TestGetRecentSensorMessages(t *testing.T) {
	sensorDataByte := []byte(`{
		"Temperature": 2,
		"Humidity": 3,
	}`)
	memoryDataByte := []byte(`{
		"Available": 1,
		"Total": 1,
		"Used": 1,
	}`)
	byteData  := [][]byte{memoryDataByte, sensorDataByte}

	virtualSensorData, memorySensorData := MapMessageData(byteData)

	expectedSensorData := &SensorData{
		Temperature: 2,
		Humidity: 3,
	}
	expectedMemoryData := &MemSensorData{
		Available: 1,
		Total: 1,
		Used: 1,
	}
	
	if virtualSensorData[0] != *expectedSensorData {
		t.Errorf("Expected sensor data to be %v but got %v", expectedSensorData, virtualSensorData[0])
	}

	if memorySensorData[0] != *expectedMemoryData {
		t.Errorf("Expected memory data to be %v but got %v", expectedMemoryData, memorySensorData[0])
	}

}