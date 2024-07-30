package app

import (
	"logbook/models"
	"logbook/models/columns"
)

type CreateObjectiveAction struct {
	Parent  models.Ovid
	Content string
	Creator columns.UserId
}
