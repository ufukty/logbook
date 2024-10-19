package augmented

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models/columns"
)

type Objective struct {
	q *database.Queries

	oid columns.ObjectiveId

	// see [Objective.fetch]
	active columns.VersionId
	obj    *database.Objective
	props  *database.Props
}

func NewObjective(q *database.Queries, oid columns.ObjectiveId) *Objective {
	return &Objective{
		q:   q,
		oid: oid,
	}
}

func (o *Objective) fetch(ctx context.Context) error {
	if o.props != nil {
		return nil
	}

	active, err := o.q.SelectActive(ctx, o.oid)
	if err != nil {
		return fmt.Errorf("oneshot.SelectActive: %w", err)
	}
	o.active = active.Vid

	obj, err := o.q.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: o.oid,
		Vid: o.active,
	})
	if err != nil {
		return fmt.Errorf("oneshot.SelectObjective: %w", err)
	}
	o.obj = &obj

	props, err := o.q.SelectProperties(ctx, obj.Pid)
	if err != nil {
		return fmt.Errorf("oneshot.SelectProperties: %w", err)
	}
	o.props = &props

	return nil
}

func (o *Objective) IsOwner(ctx context.Context, uid columns.UserId) (bool, error) {
	if err := o.fetch(ctx); err != nil {
		return false, fmt.Errorf("fetch: %w", err)
	}

	return o.props.Owner == uid, nil
}

