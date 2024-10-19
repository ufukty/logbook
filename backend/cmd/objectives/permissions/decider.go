package permissions

import (
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

type Decider struct {
	oneshot *database.Queries
	l       *logger.Logger
}

func New(pool *pgxpool.Pool, l *logger.Logger) *Decider {
	return &Decider{
		oneshot: database.New(pool),
		l:       l.Sub("decider"),
	}
}
