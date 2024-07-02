package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Persists/fcproto/internal/cloud/database"
	"github.com/Persists/fcproto/internal/cloud/database/models/entities"
	"github.com/Persists/fcproto/internal/shared/models"
	"log"
)

var ctx = context.Background()

func InsertSensorMessage(data models.SensorMessage, db *database.DB, client *entities.ClientEntity) {
	formattedContent, err := formatContent(data.Content)
	if err != nil {
		log.Printf("Failed to format content: %v", err)
		return
	}

	writtenDB, err := db.NewInsert().
		Model(&entities.SensorMessageEntity{
			Content:      formattedContent,
			ClientIpAddr: client.IpAddr,
		}).Exec(ctx)

	fmt.Println("Written to DB: ", &writtenDB)

	log.Printf("Received: Timestamp: %d, Content:\n%s", data.Timestamp, formattedContent)
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
