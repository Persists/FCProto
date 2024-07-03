package sensors

import "time"

type BaseSensor interface {
	GenerateData() string
}

// GenerateDataAtInterval generates data at specified intervals and sends it to a channel.
func GenerateDataAtInterval(vs BaseSensor, interval time.Duration, stopChan <-chan bool, dataChan chan<- string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			dataChan <- vs.GenerateData()
		case <-stopChan:
			return
		}
	}
}
