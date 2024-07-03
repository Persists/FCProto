package models

import "github.com/fatih/structs"

type MessageTopic string

const (
	All       MessageTopic = "all"
	Heartbeat MessageTopic = "heartbeat"
	Sensor    MessageTopic = "sensor"
)

type Message struct {
	Topic   MessageTopic            `json:"topic"`
	Payload *map[string]interface{} `json:"payload"`
}

func NewMessage(topic MessageTopic, payload any) Message {
	formatedPayload := structs.Map(payload)

	return Message{
		Topic:   topic,
		Payload: &formatedPayload,
	}
}
