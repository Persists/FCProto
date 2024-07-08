package messages_utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Persists/fcproto/internal/cloud/database"
	"github.com/Persists/fcproto/internal/cloud/database/models/entities"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/mitchellh/mapstructure"
)

var ctx = context.Background()

func InsertSensorMessage(db *database.DB, payload *map[string]interface{}, remoteAddress string) (err error) {
	var sensorData models.SensorMessage
	err = mapstructure.Decode(payload, &sensorData)
	if err != nil {
		log.Printf("Failed to decode sensor message: %v", err)
		return
	}

	formattedContent, err := formatContent(sensorData.Content)
	if err != nil {
		log.Printf("Failed to format content: %v", err)
		return
	}
	formattedTimestamp := time.Unix(sensorData.Timestamp, 0)

	_, err = db.NewInsert().
		Model(&entities.SensorMessageEntity{
			Content:      formattedContent,
			ClientIpAddr: remoteAddress,
			Timestamp:    formattedTimestamp,
		}).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert sensor message into database: %v", err)
		return
	}
	return nil
}

func formatContent(content string) (string, error) {
	var prettyContent interface{}
	err := json.Unmarshal([]byte(content), &prettyContent)
	if err != nil {
		// If the content is not a valid JSON, return it as is
		return content, nil
	}

	prettyJSON, err := json.MarshalIndent(prettyContent, "", "  ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}
