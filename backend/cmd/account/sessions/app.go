package sessions

import (
	"logbook/cmd/account/database"
	"logbook/internal/logger"
)

type Sessions struct {
	q *database.Queries
	l *logger.Logger
}

func New(q *database.Queries, l *logger.Logger) *Sessions {
	return &Sessions{
		q: q,
		l: l.Sub("app"),
	}
}
