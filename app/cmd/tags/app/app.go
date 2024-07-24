package app

import (
	"logbook/cmd/tags/database"
	"logbook/config/api"
	"logbook/internal/web/registryfile"
)

type App struct {
	db         *database.Queries
	internalsd *registryfile.FileReader
}

func New(db *database.Queries, apicfg *api.Config, internalsd *registryfile.FileReader) *App {
	return &App{
		db:         db,
		internalsd: internalsd,
	}
}
