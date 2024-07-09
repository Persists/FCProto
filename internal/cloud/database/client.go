package database

import (
	"fmt"
	"github.com/Persists/fcproto/internal/shared/utils"
	"log"

	serverConfig "github.com/Persists/fcproto/internal/cloud/config"
)

// DBClient manages the lifecycle of database operations
type DBClient struct {
	db *DB
}

func NewClient() *DBClient {
	return &DBClient{}
}

func (dbc *DBClient) Init(config *serverConfig.ServerConfig) {
	dbc.db = Connect(config.PostgresEnv)
}

func (dbc *DBClient) Start() error {

	_, err := dbc.db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = dbc.db.createSchema()
	if err != nil {
		fmt.Printf("Error to create Schemas")
		return nil
	}

	log.Println(utils.Colorize(utils.Green, "Database client started"))

	return nil
}

func (dbc *DBClient) Stop() error {
	err := dbc.db.Close()
	if err != nil {
		fmt.Printf("Unable to close connection to database: %v\n", err)
		return err
	}
	return nil
}

func (dbc *DBClient) GetDB() DB {
	return *dbc.db
}
