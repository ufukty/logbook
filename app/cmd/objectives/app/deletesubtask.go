package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/models"
	"logbook/models/columns"
)

type DeleteSubtaskParams struct {
	Parent, Subject models.Ovid
	Actor           columns.UserId
}

// TODO: per-viewer delta values
func (a *App) DeleteSubtask(ctx context.Context, params DeleteSubtaskParams) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := queries.New(tx)

	ap, err := a.listActivePathToRock(ctx, q, params.Subject)
	if err != nil {
		return fmt.Errorf("listActivePathToRock: %w", err)
	}

	subject, err := q.SelectObjective(ctx, queries.SelectObjectiveParams{
		Oid: params.Subject.Oid,
		Vid: params.Subject.Vid,
	})
	if err != nil {
		return fmt.Errorf("SelectObjective/subject: %w", err)
	}

	subjectprops, err := q.SelectProperties(ctx, subject.Pid)
	if err != nil {
		return fmt.Errorf("SelectProperties: %w", err)
	}

	parent, err := q.SelectObjective(ctx, queries.SelectObjectiveParams{
		Oid: params.Parent.Oid,
		Vid: params.Parent.Vid,
	})
	if err != nil {
		return fmt.Errorf("SelectObjective/parent: %w", err)
	}

	op, err := q.InsertOperation(ctx, queries.InsertOperationParams{
		Subjectoid: parent.Oid,
		Subjectvid: parent.Vid,
		Actor:      params.Actor,
		OpType:     queries.OpTypeObjDeleteSubtask,
		OpStatus:   queries.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("InsertOperation: %w", err)
	}

	_, err = q.InsertOpObjDeleteSubtask(ctx, queries.InsertOpObjDeleteSubtaskParams{
		Opid: op.Opid,
		Doid: params.Subject.Oid,
		Dvid: params.Subject.Vid,
	})
	if err != nil {
		return fmt.Errorf("InsertOpObjDeleteSubtask: %w", err)
	}

	ap[0] = models.Ovid{Oid: parent.Oid, Vid: parent.Vid}

	_, err = a.bubblink(ctx, q, ap, op, bubblinkDeltaValues{
		SubtreeCompleted: -1,
		SubtreeSize:      ternary(subjectprops.Completed, int32(-1), int32(0)),
	})
	if err != nil {
		return fmt.Errorf("bubblink: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
