package sensor

import (
	"fmt"
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	"os/exec"
	"strings"
)

// SystemInfo holds the collected system information.
type SystemInfo struct {
	Memory            types.HostMemoryInfo
	systemInformation types.HostInfo
	CPUInformation    types.CPUTimes
}

// GetSystemInfo collects and returns system information.
func GetSystemInfo() (*SystemInfo, error) {
	host, err := sysinfo.Host()
	if err != nil {
		return nil, fmt.Errorf("error getting host info: %v", err)
	}

	memInfo, err := host.Memory()
	if err != nil {
		return nil, fmt.Errorf("error getting memory info: %v", err)
	}

	cpuInfo, err := host.CPUTime()
	if err != nil {
		return nil, fmt.Errorf("error getting CPU info: %v", err)
	}

	return &SystemInfo{
		Memory:            *memInfo,
		systemInformation: host.Info(),
		CPUInformation:    cpuInfo,
	}, nil
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
