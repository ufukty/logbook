package l2

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool    *pgxpool.Pool
	oneshot *database.Queries

	l *logger.Logger
}
