package endpoints

import (
	sessions "logbook/cmd/sessions/client"
	"logbook/cmd/tags/app"
	"logbook/internal/logger"
)

type Public struct {
	a        *app.App
	l        *logger.Logger
	sessions sessions.Interface
}

func New(a *app.App, s sessions.Interface, l *logger.Logger) *Public {
	return &Public{
		a:        a,
		l:        l.Sub("endpoints/public"),
		sessions: s,
	}
}
