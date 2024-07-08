package sensors

import (
	"encoding/json"
	"log"
	"math"
)

// Analyse calculates average, standard deviation, mean, min, and max of a slice of float64
func Analyse(data []float64) (avg float64, dev float64, mean float64, min float64, max float64) {
	if len(data) == 0 {
		// Return all zeros if no data is provided
		return 0, 0, 0, 0, 0
	}

	var total float64
	var minValue float64 = math.MaxFloat64
	var maxValue float64 = -math.MaxFloat64 // Set to the smallest possible value

	for _, d := range data {
		total += d
		if d < minValue {
			minValue = d
		}
		if d > maxValue {
			maxValue = d
		}
	}

	// Calculate average
	avg = total / float64(len(data))

	// Calculate mean (which is the same as average in this context)
	mean = avg

	// Calculate standard deviation
	var sumOfSquares float64
	for _, d := range data {
		diff := d - mean
		sumOfSquares += diff * diff
	}

	dev = math.Sqrt(sumOfSquares / float64(len(data)))

	min = minValue
	max = maxValue

	return
}

// MapMessageData sorts the data into two slices, one for virtual sensor data and one for memory sensor data
func MapMessageData(data [][]byte) (virtualSensorData []SensorData, memorySensorData []MemSensorData) {
	for _, dataMsg := range data {
		var contentMap map[string]interface{}
		if err := json.Unmarshal(dataMsg, &contentMap); err != nil {
			log.Printf("Failed to unmarshal data: %v", err)
			continue
		}

		if contentMap["Temperature"] != nil {
			var sensorData SensorData
			if err := json.Unmarshal(dataMsg, &sensorData); err != nil {
				log.Printf("Failed to unmarshal data: %v", err)
				continue
			}
			virtualSensorData = append(virtualSensorData, sensorData)
		}

		if contentMap["total"] != nil {
			var memoryData MemSensorData
			if err := json.Unmarshal(dataMsg, &memoryData); err != nil {
				log.Printf("Failed to unmarshal data: %v", err)
				continue
			}
			memorySensorData = append(memorySensorData, memoryData)
		}

	}

	return
}
