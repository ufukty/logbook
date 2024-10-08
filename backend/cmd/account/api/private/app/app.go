package app

import (
	"logbook/cmd/account/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool    *pgxpool.Pool
	oneshot *database.Queries
}
