package sensors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"os/exec"
	"strings"
)

// HardwareSensor holds the collected system information.
type HardwareSensor struct {
	Memory    *mem.VirtualMemoryStat
	Load      *load.AvgStat
	CPU       []cpu.InfoStat
	IPAddress string
}

func (hs HardwareSensor) GenerateData() (dataString string) {
	transformedMem := hs.TransformedMemInfo(Megabyte)
	dataString = fmt.Sprintf(
		"Memory Total: %d MB\nMemory Available: %d MB\nMemory Used: %d MB\nLoad Average: %.2f %.2f %.2f\nIP Address: %s\n",
		transformedMem.Total,
		transformedMem.Available,
		transformedMem.Used,
		hs.Load.Load1,
		hs.Load.Load5,
		hs.Load.Load15,
		hs.IPAddress,
	)
	return
}

type MemoryUnit uint64

const (
	Megabyte MemoryUnit = 1024 * 1024
)

// GetSystemInfo collects and returns system information.
func GetSystemInfo() (*HardwareSensor, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting memory info: %v", err)
	}

	loadInfo, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("error getting load info: %v", err)
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting CPU info: %v", err)
	}

	ipAddr, err := getIPAddress()
	if err != nil {
		return nil, fmt.Errorf("error getting IP address: %v", err)
	}

	return &HardwareSensor{
		Memory:    memInfo,
		Load:      loadInfo,
		CPU:       cpuInfo,
		IPAddress: ipAddr,
	}, nil
}

func (hs HardwareSensor) TransformedMemInfo(unit MemoryUnit) *mem.VirtualMemoryStat {
	return &mem.VirtualMemoryStat{
		Total:     hs.Memory.Total / uint64(unit),
		Available: hs.Memory.Available / uint64(unit),
		Used:      hs.Memory.Used / uint64(unit),
	}
}

// Helper function to get the IP address
func getIPAddress() (string, error) {
	out, err := exec.Command("hostname", "-I").Output()
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(out))
	return ip, nil
}
