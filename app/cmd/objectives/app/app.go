package app

import (
	"logbook/cmd/objectives/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool    *pgxpool.Pool
	oneshot *queries.Queries
}

func New(pool *pgxpool.Pool) *App {
	return &App{
		pool:    pool,
		oneshot: queries.New(pool),
	}
}
