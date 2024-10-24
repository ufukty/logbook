package adapter

import (
	"context"
	"logbook/models/columns"
)

type Group struct {
}

func (a *Adapter) GetGroup(ctx context.Context, gid columns.GroupId) (*Group, error)
