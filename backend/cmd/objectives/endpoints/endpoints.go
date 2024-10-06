package endpoints

import (
	"logbook/cmd/objectives/app"
	"logbook/internal/logger"
)

type Endpoints struct {
	app *app.App
	log *logger.Logger
}

func New(app *app.App) *Endpoints {
	return &Endpoints{
		app: app,
		log: logger.New("endpoints"),
	}
}
