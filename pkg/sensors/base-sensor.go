package sensors

import (
	"encoding/json"
	"time"

	"github.com/Persists/fcproto/internal/shared/utils"
)

type BaseSensor interface {
	GenerateData() *BaseSensorData
}

type BaseSensorData struct {
	Data interface{} `json:"data"`
}

type Stringable interface {
	toString() string
}

// GenerateDataAtInterval generates data at specified intervals and sends it to a channel.
func GenerateDataAtInterval(vs BaseSensor, interval time.Duration, stopChan <-chan bool, dataChan chan<- string) {
	utils.StartTicker(interval, vs.GenerateData().ToString, stopChan, dataChan)
}

func (bsd *BaseSensorData) ToString() string {
	if bsd.Data == nil {
		return ""
	}

	dataBytes, err := json.Marshal(bsd.Data)
	if err != nil {
		return ""
	}

	return string(dataBytes)
}
