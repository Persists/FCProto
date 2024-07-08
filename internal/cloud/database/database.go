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

func Connect(env *database_models.PostgresEnv) *DB {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(env.Addr),
		pgdriver.WithDatabase(env.Database),
		pgdriver.WithUser(env.User),
		pgdriver.WithPassword(env.Password),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqlDb, pgdialect.New())

	if err := db.Ping(); err != nil {
		log.Print("Failed to connect to database")
		log.Fatal(err)
	}

	return &DB{db}
}

func (db *DB) createSchema() error {
	models := []interface{}{
		(*entities.ClientEntity)(nil),
		(*entities.SensorMessageEntity)(nil),
	}

	return db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, model := range models {
			c, err := tx.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println("schema created")
			fmt.Println(c)
		}
		return nil
	})
}

func (db *DB) InsertClient(ipAddr string) (*entities.ClientEntity, error) {
	client := &entities.ClientEntity{
		IpAddr:   ipAddr,
		LastSeen: time.Now(),
	}
	_, err := db.NewInsert().
		Model(client).
		On("CONFLICT (ip_addr) DO UPDATE").
		Set("last_seen = EXCLUDED.last_seen, notify_addr = EXCLUDED.notify_addr").
		Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert client into database: %v", err)
		return nil, err
	}
	return client, nil
}

func (db *DB) GetRecentSensorMessages() ([]entities.SensorMessageEntity, error) {
	var messages []entities.SensorMessageEntity
	tenMinutesAgo := time.Now().Add(-50 * time.Minute)

	err := db.NewSelect().
		Model(&messages).
		Where("timestamp > ?", tenMinutesAgo).
		Relation("Client").
		Order("timestamp DESC").
		Scan(ctx)
	if err != nil {
		log.Printf("Failed to get recent sensor messages: %v", err)
		return nil, err
	}

	return messages, nil
}
