package endpoints

import (
	"logbook/cmd/registration/app"
	"logbook/internal/logger"
)

type Public struct {
	a *app.App
	l *logger.Logger
}

func NewPublic(a *app.App, l *logger.Logger) *Public {
	return &Public{
		a: a,
		l: l.Sub("endpoints/public"),
	}
}
