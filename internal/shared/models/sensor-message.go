package models

import "time"

type SensorMessage struct {
	// timestamp of the message
	Timestamp int64 `json:"timestamp"`

	// message content
	Content string `json:"content"`
}

// creates a new sensor message with the given content
func NewSensorMessage(content string) SensorMessage {
	return SensorMessage{
		Timestamp: time.Now().Unix(),
		Content:   content,
	}
}
