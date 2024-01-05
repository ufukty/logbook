package endpoints

import (
	"logbook/cmd/tasks/app"
	"logbook/internal/web/logger"
)

type Endpoints struct {
	app *app.App
	log *logger.Logger
}

func NewManager(app *app.App) *Endpoints {
	return &Endpoints{
		app: app,
		log: logger.NewLogger("endpoints"),
	}
}
