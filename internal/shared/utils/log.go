package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/Persists/fcproto/internal/shared/queue"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	white  = "\033[37m"
)

const (
	Red LogColor = iota
	Green
	Yellow
	Blue
	Purple
	Cyan
	White
)

type LogColor int

func Colorize(color LogColor, message string) string {
	switch color {
	case Red:
		return fmt.Sprintf("%s%s%s", red, message, reset)
	case Green:
		return fmt.Sprintf("%s%s%s", green, message, reset)
	case Yellow:
		return fmt.Sprintf("%s%s%s", yellow, message, reset)
	case Blue:
		return fmt.Sprintf("%s%s%s", blue, message, reset)
	case Purple:
		return fmt.Sprintf("%s%s%s", purple, message, reset)
	case Cyan:
		return fmt.Sprintf("%s%s%s", cyan, message, reset)
	case White:
		return fmt.Sprintf("%s%s%s", white, message, reset)
	default:
		return message
	}
}

func LogQueuePeriodically(queue *queue.Queue[models.Message], stop chan struct{}) {
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				message := fmt.Sprintf("Messages queued for subsequent delivery: %d", queue.Len())
				log.Println(Colorize(Yellow, message))
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func FormatCloudAnalysisData(msg models.Message) (formattedText string) {
	data := *msg.Payload
	memorySensorData := data["MemorySensorData"].(map[string]interface{})
	virtualSensorData := data["VirtualSensorData"].(map[string]interface{})

	formattedText = fmt.Sprintf("The cloud recently analyzed data at %v transmitted from this device to the server, and this is the outcome of the statistical analysis:\n", msg.Time)

	formattedText += "\nThe memory sensor data transmitted from this device to the cloud is:\n"
	formattedText += fmt.Sprintf(" - Average: %e\n", memorySensorData["Average"])
	formattedText += fmt.Sprintf(" - Deviation: %.2f\n", memorySensorData["Deviation"])
	formattedText += fmt.Sprintf(" - Max: %e\n", memorySensorData["Max"])
	formattedText += fmt.Sprintf(" - Mean: %e\n", memorySensorData["Mean"])
	formattedText += fmt.Sprintf(" - Min: %e\n", memorySensorData["Min"])

	formattedText += "\nThe virtual sensor data transmitted from this device to the cloud is:\n"
	formattedText += fmt.Sprintf(" - Average: %.2f\n", virtualSensorData["Average"])
	formattedText += fmt.Sprintf(" - Deviation: %.2f\n", virtualSensorData["Deviation"])
	formattedText += fmt.Sprintf(" - Max: %.2f\n", virtualSensorData["Max"])
	formattedText += fmt.Sprintf(" - Mean: %.2f\n", virtualSensorData["Mean"])
	formattedText += fmt.Sprintf(" - Min: %.2f\n", virtualSensorData["Min"])

	return
}
