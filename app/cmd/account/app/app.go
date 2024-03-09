package app

import (
	"logbook/cmd/account/database"
)

type App struct {
	queries *database.Queries
}

func New(queries *database.Queries) *App {
	return &App{
		queries: queries,
	}
}
