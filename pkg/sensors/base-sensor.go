package sensors

import (
	"time"

	"github.com/Persists/fcproto/internal/shared/utils"
)

type BaseSensor interface {
	GenerateData() string
}

// GenerateDataAtInterval generates data at specified intervals and sends it to a channel.
func GenerateDataAtInterval(vs BaseSensor, interval time.Duration, stopChan <-chan bool, dataChan chan<- string) {
	utils.StartTicker(interval, vs.GenerateData, stopChan, dataChan)
}
