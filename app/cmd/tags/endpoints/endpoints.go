package endpoints

import (
	"logbook/cmd/tasks/database"
	"logbook/internal/web/logger"
)

type Endpoints struct {
	db  *database.Database
	log *logger.Logger
}

func NewManager(db *database.Database) *Endpoints {
	return &Endpoints{
		db:  db,
		log: logger.NewLogger("endpoints"),
	}
}
