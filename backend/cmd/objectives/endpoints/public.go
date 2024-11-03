package endpoints

import (
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/permissions"
	"logbook/internal/logger"
	"logbook/internal/rates"
)

type Public struct {
	a *app.App
	l *logger.Logger
	r *rates.Layered
	p *permissions.Decider
}

func NewPublic(a *app.App, l *logger.Logger) *Public {
	return &Public{
		a: a,
		l: l.Sub("endpoints/public"),
		r: rates.NewLayered(rates.LimiterParams{}, rates.LimiterParams{}),
	}
}
