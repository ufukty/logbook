package app

import (
	"logbook/cmd/account/app/authz"
	"logbook/cmd/account/database"
)

type App struct {
	authz   *authz.Authorization
	queries *database.Queries
}

func New(queries *database.Queries) *App {
	return &App{
		authz:   authz.New(queries),
		queries: queries,
	}
}
