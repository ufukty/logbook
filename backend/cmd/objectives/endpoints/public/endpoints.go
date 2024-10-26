package public

import (
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/permissions"
	"logbook/internal/logger"
	"logbook/internal/rates"
)

type Endpoints struct {
	a *app.App
	l *logger.Logger
	r *rates.Layered
	p *permissions.Decider
}

func New(a *app.App, l *logger.Logger) *Endpoints {
	return &Endpoints{
		a: a,
		l: l.Sub("Endpoints"),
		r: rates.NewLayered(rates.LimiterParams{}, rates.LimiterParams{}),
	}
}
