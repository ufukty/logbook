package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(connection string) (*Database, error) {
	db := &Database{}
	var err error
	db.pool, err = pgxpool.New(context.Background(), connection)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.Connect: %w", err)
	}
	return db, nil
}

func (db *Database) Close() {
	db.pool.Close()
}
