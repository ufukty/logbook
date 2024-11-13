package endpoints

import (
	"logbook/cmd/objectives/app"
	sessions "logbook/cmd/sessions/client"
	"logbook/internal/logger"
	"logbook/internal/rates"
)

type Public struct {
	a        *app.App
	l        *logger.Logger
	sessions sessions.Interface
	r        *rates.Layered
}

func NewPublic(a *app.App, sessions sessions.Interface, l *logger.Logger) *Public {
	return &Public{
		a:        a,
		sessions: sessions,
		l:        l.Sub("endpoints/public"),
		r:        rates.NewLayered(rates.LimiterParams{}, rates.LimiterParams{}),
	}
}
