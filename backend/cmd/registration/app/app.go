package app

import (
	objectives "logbook/cmd/objectives/client"
	profiles "logbook/cmd/profiles/client"
	sessions "logbook/cmd/sessions/client"
	users "logbook/cmd/users/client"
)

type App struct {
	Objectives objectives.Interface
	Profiles   profiles.Interface
	Sessions   sessions.Interface
	Users      users.Interface
}
