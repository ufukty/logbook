package permissions

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

func (d *Decider) CanGroupReadObjective(ctx context.Context, gid columns.GroupId, oid columns.ObjectiveId) error {

	return ErrUnauthorized
}

func (d *Decider) CanGroupWriteObjective(ctx context.Context, gid columns.GroupId, oid columns.ObjectiveId) error {

	return ErrUnauthorized
}

func (d *Decider) CanGroupUpdateObjective(ctx context.Context, gid columns.GroupId, oid columns.ObjectiveId) error {

	return ErrUnauthorized
}

func (d *Decider) CanGroupDeleteObjective(ctx context.Context, gid columns.GroupId, oid columns.ObjectiveId) error {

	return ErrUnauthorized
}

func (d *Decider) CanUserUpdateObjective(ctx context.Context, uid columns.UserId, oid columns.ObjectiveId) error {
	obj, err := d.db.Objective(ctx, oid)
	if err != nil {
		return fmt.Errorf("db.Objective: %w", err)
	}

	ok, err := obj.IsOwner(ctx, uid)
	if err != nil {
		return fmt.Errorf("obj.IsOwner: %w", err)
	} else if ok {
		return nil
	}

	// Collaborator list

	// Check membership to groups that are assigned as collaborator to the objective

	// Recur on parent (inheritance)

	return ErrUnauthorized
}
