package endpoints

import (
	"logbook/cmd/profiles/app"
	"logbook/internal/logger"
)

type Private struct {
	a *app.App
	l *logger.Logger
}

func NewPrivate(a *app.App, l *logger.Logger) *Private {
	return &Private{
		a: a,
		l: l.Sub("endpoints/public"),
	}
}
