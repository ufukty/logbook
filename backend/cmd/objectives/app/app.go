package app

import (
	"logbook/cmd/objectives/app/l2"
	"logbook/cmd/objectives/database"
	"logbook/internal/logger"
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

func newCacheStore(l *logger.Logger) *caches {
	return &caches{
		Uss: stores.NewFixedSizeKV[usssubject, int32](l.Sub("uss cache"), 100000), // 108byte * 100.000 = 10.8mb
	}
}

type App struct {
	pool    *pgxpool.Pool
	oneshot *database.Queries
	caches  *caches
	l2      *l2.App

	l *logger.Logger
}

func New(pool *pgxpool.Pool, l *logger.Logger) *App {
	l = l.Sub("App")
	return &App{
		pool:    pool,
		oneshot: database.New(pool),
		caches:  newCacheStore(l),
		l:       l,
	}
}
