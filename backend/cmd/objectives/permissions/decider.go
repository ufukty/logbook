package permissions

import (
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/permissions/adapter"
	"logbook/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

type Decider struct {
	db *adapter.Adapter
	l  *logger.Logger
}

func New(pool *pgxpool.Pool, l *logger.Logger) *Decider {
	return &Decider{
		db: adapter.New(database.New(pool)),
		l:  l.Sub("decider"),
	}
}
