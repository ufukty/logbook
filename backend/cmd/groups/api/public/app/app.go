package app

import (
	"logbook/cmd/groups/database"
	"logbook/models"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5/pgxpool"
)

type usssubject struct {
	Viewer columns.UserId
	Object models.Ovid
}

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
