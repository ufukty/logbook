package endpoints

import (
	"logbook/cmd/account/app"
	"logbook/internal/logger"
)

type Endpoints struct {
	app *app.App
	l   logger.Logger
}

func New(a *app.App) *Endpoints {
	return &Endpoints{
		app: a,
		l:   *logger.New("Endpoints"),
	}
}
