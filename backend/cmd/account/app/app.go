package app

import (
	"logbook/cmd/account/database"
	"logbook/cmd/account/permissions"
	"logbook/cmd/account/sessions"
	objectives "logbook/cmd/objectives/client"
	"logbook/config/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pd      *permissions.Decider
	pool    *pgxpool.Pool
	oneshot *database.Queries
	s       *sessions.Sessions

	objectives *objectives.Client
}

func New(pool *pgxpool.Pool, apicfg *api.Config, objectivesctl *objectives.Client) *App {
	return &App{
		pd:         permissions.New(database.New(pool)),
		pool:       pool,
		oneshot:    database.New(pool),
		objectives: objectivesctl,
	}
}
