package database

import (
	"context"
	"testing"

	serverConfig "github.com/Persists/fcproto/internal/cloud/config"
	"github.com/stretchr/testify/assert"
)

func TestDBClient_Integration(t *testing.T) {
	// Initialize a server configuration for testing
	config := &serverConfig.ServerConfig{
		PostgresEnv: "your_test_db_connection_string",
	}

	client := NewClient()
	err := client.Init(config)
	assert.NoError(t, err, "Expected no error during Init")

	err = client.Start()
	assert.NoError(t, err, "Expected no error during Start")

	db := client.GetDB()

	// Example integration test, for instance creating a table
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test_table (id UUID PRIMARY KEY, name TEXT);`)
	assert.NoError(t, err, "Expected no error during table creation")

	// Clean up
	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS test_table;`)
	assert.NoError(t, err, "Expected no error during table drop")

	err = client.Stop()
	assert.NoError(t, err, "Expected no error during Stop")
}