package adapter

import "logbook/cmd/objectives/database"

type Adapter struct {
	q *database.Queries
}

func New(q *database.Queries) *Adapter {
	return &Adapter{
		q: q,
	}
}
