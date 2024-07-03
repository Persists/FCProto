package models

import "time"

type SensorMessage struct {
	// timestamp of the message
	Timestamp int64 `json:"timestamp"`

	// message content
	// TODO: expand the content to include more fields
	Content string `json:"content"`
}

func NewSensorMessage(content string) SensorMessage {
	return SensorMessage{
		Timestamp: time.Now().Unix(),
		Content:   content,
	}
}
