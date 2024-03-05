package endpoints

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/web/logger"
)

type Endpoints struct {
	db  *database.Queries
	log *logger.Logger
}

func NewManager(db *database.Queries) *Endpoints {
	return &Endpoints{
		db:  db,
		log: logger.NewLogger("endpoints"),
	}
}
