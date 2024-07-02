package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Persists/fcproto/internal/server/database/models"
	"github.com/Persists/fcproto/internal/server/database/models/entities"
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
		pgdriver.WithDatabase(env.Database),
		pgdriver.WithUser(env.User),
		pgdriver.WithPassword(env.Password),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqlDb, pgdialect.New())

	return &DB{db}
}

func (db *DB) CloseConnection() {
	err := db.Close()
	if err != nil {
		fmt.Printf("Unable to close connection to database: %v\n", err)
	}
}

func (db *DB) Start() {
	_, err := db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		fmt.Println(err)
	}

	err = db.createSchema()
	if err != nil {
		fmt.Printf("Error to create Schemas")
	}
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
			fmt.Println(c)
		}
		return nil
	})
}
