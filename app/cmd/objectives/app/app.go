package app

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/web/balancer"
)

type App struct {
	queries  *database.Queries
	internal balancer.InstanceSource
}

func New(queries *database.Queries, internal balancer.InstanceSource) *App {
	return &App{
		queries:  queries,
		internal: internal,
	}
}
