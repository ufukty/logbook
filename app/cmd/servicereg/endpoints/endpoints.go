package endpoints

import (
	"logbook/cmd/servicereg/app"
)

type Endpoints struct {
	a *app.App
}

func New(a *app.App) *Endpoints {
	return &Endpoints{
		a: a,
	}
}
