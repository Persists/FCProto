package models

import (
	"time"

	"github.com/fatih/structs"
)

type MessageTopic string

const (
	All       MessageTopic = "all"
	Heartbeat MessageTopic = "heartbeat"
	Sensor    MessageTopic = "sensor"
	Fog       MessageTopic = "edge"
	Cloud     MessageTopic = "cloud"
	Analysis  MessageTopic = "Analysis"
)

type Message struct {
	Topic   MessageTopic            `json:"topic"`
	Payload *map[string]interface{} `json:"payload"`
	Time    time.Time               `json:"time"`
}

// creates a new message with the given topic and payload
func NewMessage(topic MessageTopic, payload any) Message {
	formatedPayload := structs.Map(payload)

	return Message{
		Topic:   topic,
		Payload: &formatedPayload,
		Time:    time.Now(),
	}
}
