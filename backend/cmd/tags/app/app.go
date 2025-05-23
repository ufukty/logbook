package app

import (
	"logbook/cmd/tags/database"
	"logbook/internal/web/registryfile"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool       *pgxpool.Pool
	oneshot    *database.Queries
	internalsd *registryfile.FileReader
}

func New(pool *pgxpool.Pool, internalsd *registryfile.FileReader) *App {
	return &App{
		pool:       pool,
		oneshot:    database.New(pool),
		internalsd: internalsd,
	}
}
