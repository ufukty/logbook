package endpoints

import "logbook/cmd/account/app"

type Endpoints struct {
	app *app.App
}

func New(a *app.App) *Endpoints {
	return &Endpoints{
		app: a,
	}
}
