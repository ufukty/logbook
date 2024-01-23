package app

import "logbook/cmd/tasks/database"

type CreateObjectiveAction struct {
	Parent  database.Ovid
	Content string
	Creator database.UserId
}
