package app

import (
	"logbook/cmd/profiles/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool    *pgxpool.Pool
	oneshot *database.Queries
}

func New(pool *pgxpool.Pool) *App {
	return &App{
		pool:    pool,
		oneshot: database.New(pool),
	}
}
