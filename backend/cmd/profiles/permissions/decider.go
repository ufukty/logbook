package permissions

import (
	"fmt"
	"logbook/cmd/profiles/database"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

type Decider struct {
	q *database.Queries
}

func New(q *database.Queries) *Decider {
	return &Decider{
		q: q,
	}
}
