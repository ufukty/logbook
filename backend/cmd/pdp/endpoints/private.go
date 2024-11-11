package endpoints

import (
	"logbook/cmd/pdp/decider"
	"logbook/internal/logger"
)

type Private struct {
	d *decider.Decider
	l *logger.Logger
}

func NewPrivate(d *decider.Decider, l *logger.Logger) *Private {
	return &Private{
		d: d,
		l: l.Sub("endpoints/private"),
	}
}
