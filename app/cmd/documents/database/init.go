package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(connection string) (*Database, error) {
	db := &Database{}
	var err error
	db.pool, err = pgxpool.Connect(context.Background(), connection)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize Database connection using pgx: %w", err)
	}
	return db, nil
}

func (db *Database) Close() {
	db.pool.Close()
}
