package database

import (
	"fmt"
	serverConfig "github.com/Persists/fcproto/internal/cloud/config"
)

// DatabaseManager manages the lifecycle of database operations
type DatabaseManager struct {
	db *DB
}

func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

func (manager *DatabaseManager) Init(config *serverConfig.ServerConfig) error {
	manager.db = Connect(config.PostgresEnv)
	return nil
}

func (manager *DatabaseManager) Start() error {
	_, err := manager.db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = manager.db.createSchema()
	if err != nil {
		fmt.Printf("Error to create Schemas")
		return nil
	}

	return nil
}

func (manager *DatabaseManager) Stop() error {
	err := manager.db.Close()
	if err != nil {
		fmt.Printf("Unable to close connection to database: %v\n", err)
		return err
	}
	return nil
}

func (manager *DatabaseManager) GetDB() DB {
	return *manager.db
}
