package models

type AnalysisVirtualSensorData struct {
	Average   float64 `json:"average"`
	Deviation float64 `json:"deviation"`
	Mean      float64 `json:"mean"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
}

type AnalysisMemorySensorData struct {
	Average   float64 `json:"average"`
	Deviation float64 `json:"deviation"`
	Mean      float64 `json:"mean"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
}

type AnalysisPayload struct {
	VirtualSensorData AnalysisVirtualSensorData `json:"virtual_sensor_data"`
	MemorySensorData  AnalysisMemorySensorData  `json:"memory_sensor_data"`
}
