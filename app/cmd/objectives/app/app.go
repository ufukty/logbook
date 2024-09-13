package app

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/stores"
	"logbook/models"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5/pgxpool"
)

type usssubject struct {
	Viewer columns.UserId
	Object models.Ovid
}

type caches struct {
	Uss *stores.FixedSizeKV[usssubject, int32]
}

func newCacheStore() *caches {
	return &caches{
		Uss: stores.NewFixedSizeKV[usssubject, int32]("uss cache", 100000), // 108byte * 100.000 = 10.8mb
	}
}

type App struct {
	pool    *pgxpool.Pool
	oneshot *database.Queries
	caches  *caches
}

func New(pool *pgxpool.Pool) *App {
	return &App{
		pool:    pool,
		oneshot: database.New(pool),
		caches:  newCacheStore(),
	}
}
