package app

import (
	"logbook/cmd/account/api/public/app/authz"
	"logbook/cmd/account/database"
	objectives "logbook/cmd/objectives/client"
	"logbook/config/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	authz   *authz.Authorization
	pool    *pgxpool.Pool
	oneshot *database.Queries

	objectives *objectives.Client
}

func New(pool *pgxpool.Pool, apicfg *api.Config, objectivesctl *objectives.Client) *App {
	return &App{
		authz:      authz.New(database.New(pool)),
		pool:       pool,
		oneshot:    database.New(pool),
		objectives: objectivesctl,
	}
}
