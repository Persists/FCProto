package database

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/Persists/fcproto/internal/cloud/database/models/entities"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/mitchellh/mapstructure"
)

// InsertSensorMessage inserts a sensor message into the database
func InsertSensorMessage(db *DB, payload *map[string]interface{}, remoteAddress string) (err error) {
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

	client := UpsertClient(db, remoteAddress)
	_, err = db.NewInsert().
		Model(&entities.SensorMessageEntity{
			Content:      formattedContent,
			Client:       client,
			ClientIpAddr: client.IpAddr,
			Timestamp:    formattedTimestamp,
		}).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert sensor message into database: %v", err)
		return
	}
	return nil
}

// UpsertClient inserts a new client into the database if it doesn't exist, otherwise updates the last seen time
func UpsertClient(db *DB, remoteAddr string) *entities.ClientEntity {
	newClient := new(entities.ClientEntity)
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Printf("Failed to split remote address: %v", err)
	}
	newClient.IpAddr = ip

	db.NewInsert().
		Model(newClient).
		On("CONFLICT (ip_addr) DO UPDATE").
		Set("ip_addr = EXCLUDED.ip_addr, last_seen = NOW()").
		Returning("*").
		Exec(ctx)

	return newClient
}

// GetRecentSensorMessages returns the most recent sensor messages (1 minute) of a client
func GetRecentSensorMessages(db *DB, ipAddr string) (sensorMessages [][]byte, insertErr error) {
	recentTime := time.Now().Add(-1 * time.Minute)

	var contentArray []string
	err := db.NewSelect().
		Model((*entities.SensorMessageEntity)(nil)).
		Where("client_ip_addr = ?", ipAddr).
		Where("timestamp > ?", recentTime).
		Order("timestamp DESC").
		Column("content").
		Scan(ctx, &contentArray)
	if err != nil {
		log.Printf("Failed to get sensor messages of last 1 minute: %v", err)
		return nil, err
	}

	for _, content := range contentArray {
		sensorMessages = append(sensorMessages, []byte(content))
	}

	return sensorMessages, nil
}

// formatContent formats the content of the sensor message
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
