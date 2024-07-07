package sensors

import (
	"github.com/shirou/gopsutil/v3/load"
)

type CpuSensor struct {
	BaseSensor
}

type CpuSensorData struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type MappedCpuSensorData struct {
	Data CpuSensorData `json:"data"`
}

func (s *CpuSensor) GenerateData() *BaseSensorData {
	avg, err := load.Avg()

	if err != nil {
		return nil
	}

	data := &CpuSensorData{
		Load1:  avg.Load1,
		Load5:  avg.Load5,
		Load15: avg.Load15,
	}

	return &BaseSensorData{Data: data}
}

func NewCpuSensor() BaseSensor {
	return &CpuSensor{}
}
