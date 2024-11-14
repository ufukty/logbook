package app

import (
	"context"
	"logbook/models/columns"
	"logbook/models/transports"
)

func (a *App) GrantEmail(ctx context.Context, email columns.Email) (transports.EmailGrant, error)

func (a *App) GrantPassword(ctx context.Context, password transports.Password) (transports.PasswordGrant, error)

func (a *App) GrantPhone(ctx context.Context, phone columns.Phone) (transports.PhoneGrant, error)
