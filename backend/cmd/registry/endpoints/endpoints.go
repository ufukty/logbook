package endpoints

import (
	"logbook/cmd/registry/app"
)

type Endpoints struct {
	a *app.App
}

func New(a *app.App) *Endpoints {
	return &Endpoints{
		a: a,
	}
}
