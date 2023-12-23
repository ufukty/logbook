package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func initDatabaseConnection(connection string) {
	var err error
	pool, err = pgxpool.Connect(context.Background(), connection)
	if err != nil {
		log.Fatalf("Could not initialize Database connection using pgx.\n^ Error details: %s", err)
	}
}

func Init(connection string) {
	initDatabaseConnection(connection)
}

func Close() {
	pool.Close()
}
