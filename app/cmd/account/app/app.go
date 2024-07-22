package app

import (
	"logbook/cmd/account/app/authz"
	"logbook/cmd/account/database"
	objectives "logbook/cmd/objectives/client"
	"logbook/config/api"
	"logbook/internal/web/balancer"
)

type App struct {
	authz      *authz.Authorization
	queries    *database.Queries
	objectives *objectives.Client
}

func New(queries *database.Queries, apicfg *api.Config, internalsd balancer.InstanceSource) *App {
	return &App{
		authz:      authz.New(queries),
		queries:    queries,
		objectives: objectives.NewClient(balancer.New(internalsd), apicfg),
	}
}
