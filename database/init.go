package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func Init(connection string) {
	var err error
	pool, err = pgxpool.Connect(context.Background(), connection)
	if err != nil {
		log.Fatalf("Could not initialize Database connection using pgx.\n^ Error details: %s", err)
	}
}

func Close() {
	pool.Close()
}
