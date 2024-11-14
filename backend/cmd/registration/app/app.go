package app

import (
	objectives "logbook/cmd/objectives/client"
	profiles "logbook/cmd/profiles/client"
	sessions "logbook/cmd/sessions/client"
	users "logbook/cmd/users/client"
	"logbook/internal/stores"
	"logbook/models/columns"
	"logbook/models/transports"
)

type grants struct {
	phone    stores.FixedSizeKV[transports.PhoneGrant, columns.Phone]
	email    stores.FixedSizeKV[transports.EmailGrant, columns.Email]
	password stores.FixedSizeKV[transports.PasswordGrant, transports.Password]
}

type App struct {
	Objectives objectives.Interface
	Profiles   profiles.Interface
	Sessions   sessions.Interface
	Users      users.Interface
	grants     grants
}
