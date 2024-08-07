package app

import (
	"logbook/cmd/objectives/database"
)

type App struct {
	queries *database.Queries
}

func New(queries *database.Queries) *App {
	return &App{
		queries: queries,
	}
}
