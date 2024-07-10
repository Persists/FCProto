package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	database_models "github.com/Persists/fcproto/internal/cloud/database/models"
	"github.com/Persists/fcproto/internal/cloud/database/models/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DB struct {
	*bun.DB
}

var ctx = context.Background()

// Connect creates a new database connection
// using exponential backoff when failing to connect
func Connect(env *database_models.PostgresEnv) *DB {
	retryDelay := 1 * time.Second
	maxRetryDelay := 16 * retryDelay

	for {
		sqlDb := sql.OpenDB(pgdriver.NewConnector(
			pgdriver.WithAddr(env.Addr),
			pgdriver.WithDatabase(env.Database),
			pgdriver.WithUser(env.User),
			pgdriver.WithPassword(env.Password),
			pgdriver.WithInsecure(true),
		))
		db := bun.NewDB(sqlDb, pgdialect.New())

		if err := db.Ping(); err != nil {
			log.Printf("Failed to connect to Database: %v. Retrying in %s...", err, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2
			if retryDelay > maxRetryDelay {
				retryDelay = maxRetryDelay
			}
			continue
		}

		return &DB{db}
	}
}

// createSchema creates the database schema
func (db *DB) createSchema() error {
	models := []interface{}{
		(*entities.ClientEntity)(nil),
		(*entities.SensorMessageEntity)(nil),
	}

	return db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, model := range models {
			_, err := tx.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		return nil
	})
}

// InsertClient inserts a new client into the database
func (db *DB) InsertClient(ipAddr string) (*entities.ClientEntity, error) {
	client := &entities.ClientEntity{
		IpAddr:   ipAddr,
		LastSeen: time.Now(),
	}
	_, err := db.NewInsert().
		Model(client).
		On("CONFLICT (ip_addr) DO UPDATE").
		Set("last_seen = EXCLUDED.last_seen").
		Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert client into database: %v", err)
		return nil, err
	}
	return client, nil
}
