package app

import db "logbook/cmd/tasks/database"

type App struct {
	db *db.Database
}

func New(db *db.Database) *App {
	return &App{db: db}
}
