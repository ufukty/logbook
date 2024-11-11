package endpoints

import (
	"logbook/cmd/objectives/app"
	"logbook/internal/logger"
	"logbook/internal/rates"
)

type Public struct {
	a *app.App
	l *logger.Logger
	r *rates.Layered
}

func NewPublic(a *app.App, l *logger.Logger) *Public {
	return &Public{
		a: a,
		l: l.Sub("endpoints/public"),
		r: rates.NewLayered(rates.LimiterParams{}, rates.LimiterParams{}),
	}
}
