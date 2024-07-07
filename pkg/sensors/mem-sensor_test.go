package sensors

import "testing"

func TestGenerateMemoryData(t *testing.T) {
	mem := NewMemSensor()

	data := mem.GenerateData()

	if data == nil || data.Data == nil {
		t.Errorf("Expected data but got nil")
	}

	memData := data.Data.(*MemSensorData)

	if memData.Available > memData.Total {
		t.Errorf("Expected available memory between 0 and %d but got %d", memData.Total, memData.Available)
	}

	if memData.Used > memData.Total {
		t.Errorf("Expected used memory between 0 and %d but got %d", memData.Total, memData.Used)
	}

}
