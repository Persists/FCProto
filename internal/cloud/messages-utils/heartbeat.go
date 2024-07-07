package messages_utils

import (
	"fmt"
	"github.com/Persists/fcproto/internal/cloud/database"
	"github.com/Persists/fcproto/internal/cloud/database/models/entities"
	"github.com/Persists/fcproto/internal/shared/models"
	"github.com/mitchellh/mapstructure"
	"log"
	"net"
	"strings"
)

// var startCounter = 0

func UpdateLastSeen(db *database.DB, client *entities.ClientEntity) {
	_, err := db.NewUpdate().
		Model(&entities.ClientEntity{}).
		Where("ip_addr = ?", client.IpAddr).
		Set("last_seen = NOW()").
		Exec(ctx)

	if err != nil {
		log.Printf("Failed to update last seen: %v", err)
	}
}

/*
func UpdateLastSeen(db *database.DB, payload *map[string]interface{}, client *entities.ClientEntity) {
	if startCounter == 0 {
		updateCallbackAddr(db, payload, client)
		startCounter++
	}

	_, err := db.NewUpdate().
		Model(&entities.ClientEntity{}).
		Where("ip_addr = ?", client.IpAddr).
		Set("last_seen = NOW()").
		Exec(ctx)

	if err != nil {
		log.Printf("Failed to update last seen: %v", err)
	}
}*/

func UpdateNotifyAddr(db *database.DB, payload *map[string]interface{}, client *entities.ClientEntity) error {
	var heartbeatData models.HeartbeatMessage
	err := mapstructure.Decode(payload, &heartbeatData)
	if err != nil {
		log.Printf("Failed to decode sensor message: %v", err)
		return err
	}
	ip := strings.Split(client.IpAddr, ":")[0]
	heartbeatData.CallbackPort = fmt.Sprintf("%s:%s", ip, heartbeatData.CallbackPort)

	var currentCallbackAddr string
	err = db.NewSelect().
		Model(&entities.ClientEntity{}).
		Column("notify_addr").
		Where("ip_addr = ?", client.IpAddr).
		Scan(ctx, &currentCallbackAddr)
	if err != nil {
		log.Printf("Failed to get notify addr: %v", err)
		return err
	}

	if currentCallbackAddr == heartbeatData.CallbackPort {
		return nil
	}

	_, err = db.NewUpdate().
		Model(client).
		Set("notify_addr = ?", heartbeatData.CallbackPort).
		Where("ip_addr = ?", client.IpAddr).
		Exec(ctx)

	if err != nil {
		log.Printf("Failed to update notify addr: %v", err)
		return err
	}

	fmt.Printf("Updated notify addr to: %s\n", heartbeatData.CallbackPort)
	return nil
}

func NotifyAllClients(db *database.DB) {
	var clients []entities.ClientEntity

	err := db.NewSelect().
		DistinctOn("notify_addr").
		Model(&clients).
		Column("notify_addr").
		Where("notify_addr IS NOT NULL").
		Scan(ctx)
	if err != nil {
		log.Printf("Failed to get callback addr: %v", err)
		return
	}

	var (
		notifiedClients    []entities.ClientEntity
		unreachableClients []entities.ClientEntity
	)

	for _, client := range clients {
		conn, err := net.Dial("tcp", client.NotifyAddr)
		if err != nil {
			unreachableClients = append(unreachableClients, client)
			continue
		}
		conn.Close()
		notifiedClients = append(notifiedClients, client)
	}
	fmt.Printf("In total %d notified clients: %v\n", len(notifiedClients), notifiedClients)
	fmt.Printf("In total %d unreachable clients: %v\n", len(unreachableClients), unreachableClients)
}
