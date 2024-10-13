package endpoints

import (
	"logbook/cmd/objectives/app"
	"logbook/internal/logger"
	"logbook/internal/rates"
)

type Endpoints struct {
	a *app.App
	l *logger.Logger
	r *rates.Layered
}

func New(a *app.App, l *logger.Logger) *Endpoints {
	return &Endpoints{
		a: a,
		l: l.Sub("Endpoints"),
		r: rates.NewLayered(rates.LimiterParams{}, rates.LimiterParams{}),
	}
}
