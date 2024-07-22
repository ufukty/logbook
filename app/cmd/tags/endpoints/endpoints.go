package endpoints

import (
	"logbook/cmd/tags/app"
	"logbook/internal/web/logger"
)

type Endpoints struct {
	app *app.App
	log *logger.Logger
}

func New(app *app.App) *Endpoints {
	return &Endpoints{
		app: app,
		log: logger.NewLogger("endpoints"),
	}
}
