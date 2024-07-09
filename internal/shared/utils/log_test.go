package utils

import (
	"fmt"
	"github.com/Persists/fcproto/internal/shared/models"
	"strings"
	"testing"
)

func TestFormatCloudAnalysisData(t *testing.T) {
	payload := models.AnalysisPayload{
		VirtualSensorData: models.AnalysisVirtualSensorData{
			Average:   10.60869565476542,
			Deviation: 12.739220970376122,
			Mean:      10.60869565476542,
			Min:       -9.300000190734863,
			Max:       29.299999237060547,
		},
		MemorySensorData: models.AnalysisMemorySensorData{
			Average:   6.7171086336e+10,
			Deviation: 0,
			Mean:      6.7171086336e+10,
			Min:       6.7171086336e+10,
			Max:       6.7171086336e+10,
		},
	}

	msg := models.NewMessage(models.All, payload)
	timestamp := msg.Time

	expectedOutput := fmt.Sprintf(`The cloud recently analyzed data at %v transmitted from this device to the server, and this is the outcome of the statistical analysis:

The memory sensor data transmitted from this device to the cloud is:
 - Average: 6.717109e+10
 - Deviation: 0.00
 - Max: 6.717109e+10
 - Mean: 6.717109e+10
 - Min: 6.717109e+10

The virtual sensor data transmitted from this device to the cloud is:
 - Average: 10.61
 - Deviation: 12.74
 - Max: 29.30
 - Mean: 10.61
 - Min: -9.30

`, timestamp)

	result := FormatCloudAnalysisData(msg)

	if strings.TrimSpace(result) != strings.TrimSpace(expectedOutput) {
		t.Errorf("test failed. Expected: %s, got: %s", expectedOutput, result)
	}
}
