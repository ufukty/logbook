package endpoints

import (
	"logbook/cmd/groups/api/public/app"
	"logbook/internal/logger"
)

type Endpoints struct {
	a *app.App
	l *logger.Logger
}

func New(a *app.App, l *logger.Logger) *Endpoints {
	return &Endpoints{
		a: a,
		l: l.Sub("Endpoints"),
	}
}
