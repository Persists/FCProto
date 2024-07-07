package sensors

import (
	"testing"
)

func TestGenerateCPUData(t *testing.T) {
	cpu := NewCpuSensor()

	data := cpu.GenerateData()

	if data == nil || data.Data == nil {
		t.Errorf("Expected data but got nil")
	}

	cpuData := data.Data.(*CpuSensorData)

	if cpuData.Load1 < 0 || cpuData.Load1 > 100 {
		t.Errorf("Expected load1 between 0 and 100 but got %f", cpuData.Load1)
	}

	if cpuData.Load5 < 0 || cpuData.Load5 > 100 {
		t.Errorf("Expected load5 between 0 and 100 but got %f", cpuData.Load5)
	}

	if cpuData.Load15 < 0 || cpuData.Load15 > 100 {
		t.Errorf("Expected load15 between 0 and 100 but got %f", cpuData.Load15)
	}

}
