package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func Init() {
	urlExample := "postgres://postgres:password@localhost:5432/testdatabase" // os.Getenv("DATABASE_URL")
	var errPgxPool error
	pool, errPgxPool = pgxpool.Connect(context.Background(), urlExample)
	if errPgxPool != nil {
		log.Fatalf("Could not initialize Database connection using pgx %s", errPgxPool)
	}
}

func Close() {
	pool.Close()
}
