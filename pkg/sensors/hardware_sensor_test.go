package sensors

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func TestHardwareSensor_GenerateData(t *testing.T) {
	hs := HardwareSensor{
		Memory: &mem.VirtualMemoryStat{
			Total:     8192 * 1024 * 1024,
			Available: 4096 * 1024 * 1024,
			Used:      4096 * 1024 * 1024,
		},
		Load: &load.AvgStat{
			Load1:  0.5,
			Load5:  1.0,
			Load15: 1.5,
		},
		CPU: []cpu.InfoStat{
			{
				CPU:      0,
				VendorID: "GenuineIntel",
				Family:   "",
				Model:    "",
				Cores:    4,
				Mhz:      2000,
			},
		},
		IPAddress: "192.168.1.100",
	}

	dataString := hs.GenerateData()

	data := &HardwareSensorData{}
	err := json.Unmarshal([]byte(dataString), data)

	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}

	if data.MemoryTotal != hs.Memory.Total/1024/1024 {
		t.Errorf("Expected total memory to be %d but got %d", hs.Memory.Total/1024/1024, data.MemoryTotal)
	}

	if data.MemoryAvailable != hs.Memory.Available/1024/1024 {
		t.Errorf("Expected available memory to be %d but got %d", hs.Memory.Available/1024/1024, data.MemoryAvailable)
	}

	if data.MemoryUsed != hs.Memory.Used/1024/1024 {
		t.Errorf("Expected used memory to be %d but got %d", hs.Memory.Used/1024/1024, data.MemoryUsed)
	}

	if data.LoadAverage != fmt.Sprintf("%.2f %.2f %.2f", hs.Load.Load1, hs.Load.Load5, hs.Load.Load15) {
		t.Errorf("Expected load average to be %s but got %s", fmt.Sprintf("%.2f %.2f %.2f", hs.Load.Load1, hs.Load.Load5, hs.Load.Load15), data.LoadAverage)
	}

	if data.IpAddress != hs.IPAddress {
		t.Errorf("Expected IP address to be %s but got %s", hs.IPAddress, data.IpAddress)
	}

}
