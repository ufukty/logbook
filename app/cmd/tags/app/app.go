package app

import (
	"logbook/cmd/tags/database"
	"logbook/config/api"
	"logbook/internal/web/discoveryfile"
)

type App struct {
	db         *database.Queries
	internalsd *discoveryfile.FileReader
}

func New(db *database.Queries, apicfg *api.Config, internalsd *discoveryfile.FileReader) *App {
	return &App{
		db:         db,
		internalsd: internalsd,
	}
}
