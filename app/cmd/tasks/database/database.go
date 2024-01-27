package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	db *pgxpool.Pool
}

func New(connection string) (*Queries, error) {
	q := &Queries{}
	var err error
	q.db, err = pgxpool.New(context.Background(), connection)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.Connect: %w", err)
	}
	return q, nil
}

func (q *Queries) Close() {
	q.db.Close()
}

// func (q *Queries) WithTx(tx pgx.Tx) *Queries {
// 	return &Queries{
// 		db: tx,
// 	}
// }
