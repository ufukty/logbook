package adapter

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models/columns"
)

type Objective struct {
	Oid    columns.ObjectiveId
	Active columns.VersionId
	Obj    *database.Objective
	Props  *database.Props
}

func (o *Objective) IsOwner(ctx context.Context, uid columns.UserId) (bool, error) {
	return o.Props.Owner == uid, nil
}

func (o *Objective) ListCollaboratorsCanRead(ctx context.Context) ([]columns.UserId, error)

func (a *Adapter) Objective(ctx context.Context, oid columns.ObjectiveId) (*Objective, error) {
	active, err := a.q.SelectActive(ctx, oid)
	if err != nil {
		return nil, fmt.Errorf("oneshot.SelectActive: %w", err)
	}

	obj, err := a.q.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: oid,
		Vid: active.Vid,
	})
	if err != nil {
		return nil, fmt.Errorf("oneshot.SelectObjective: %w", err)
	}

	props, err := a.q.SelectProperties(ctx, obj.Pid)
	if err != nil {
		return nil, fmt.Errorf("oneshot.SelectProperties: %w", err)
	}

	return &Objective{
		Oid:    oid,
		Active: active.Vid,
		Obj:    &obj,
		Props:  &props,
	}, nil
}

