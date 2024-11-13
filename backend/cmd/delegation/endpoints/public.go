package endpoints

import (
	"logbook/cmd/groups/app"
	sessions "logbook/cmd/sessions/client"
	"logbook/internal/logger"
)

type Public struct {
	a        *app.App
	l        *logger.Logger
	sessions sessions.Interface
}

func NewPublic(a *app.App, l *logger.Logger) *Public {
	return &Public{
		a: a,
		l: l.Sub("endpoints/public"),
	}
}
