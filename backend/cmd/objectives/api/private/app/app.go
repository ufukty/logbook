package app

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool    *pgxpool.Pool
	oneshot *database.Queries

	l *logger.Logger
}

func New(pool *pgxpool.Pool, l *logger.Logger) *App {
	return &App{
		pool:    pool,
		oneshot: database.New(pool),
		l:       l.Sub("App"),
	}
}
