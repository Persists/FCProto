package sensors

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type MemSensor struct {
	BaseSensor
}

type MemSensorData struct {
	Available uint64 `json:"available"`
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
}

func (s *MemSensor) GenerateData() *BaseSensorData {
	memInfo, err := mem.VirtualMemory()

	if err != nil {
		return nil
	}

	data := &MemSensorData{
		Available: memInfo.Available,
		Total:     memInfo.Total,
		Used:      memInfo.Used,
	}

	return &BaseSensorData{Data: data}
}

func NewMemSensor() *MemSensor {
	return &MemSensor{}
}
