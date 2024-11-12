package endpoints

import (
	account "logbook/cmd/account/client"
	"logbook/cmd/groups/app"
	"logbook/internal/logger"
)

type Public struct {
	a        *app.App
	l        *logger.Logger
	accounts account.Interface
}

func NewPublic(a *app.App, l *logger.Logger) *Public {
	return &Public{
		a: a,
		l: l.Sub("endpoints/public"),
	}
}
