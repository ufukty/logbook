package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/models"
	"logbook/models/columns"
)

type DeleteSubtaskParams struct {
	Subject models.Ovid
	Actor   columns.UserId
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
	parentOvid := ap[1]

	parent, err := q.SelectObjective(ctx, queries.SelectObjectiveParams{
		Oid: parentOvid.Oid,
		Vid: parentOvid.Vid,
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

	subjectbups, err := q.SelectBottomUpProps(ctx, subject.Bupid)
	if err != nil {
		return fmt.Errorf("SelectBottomUpProps/subject: %w", err)
	}

	parentbups, err := q.SelectBottomUpProps(ctx, parent.Bupid)
	if err != nil {
		return fmt.Errorf("SelectBottomUpProps/parent: %w", err)
	}

	deltas := bubblinkDeltaValues{
		SubtreeCompleted: -1 * (subjectbups.SubtreeCompleted + ternary(subjectprops.Completed, int32(1), int32(0))),
		SubtreeSize:      -1 * (subjectbups.SubtreeSize + 1),
	}
	parentbups.Children += -1
	parentbups.SubtreeCompleted += deltas.SubtreeCompleted
	parentbups.SubtreeSize += deltas.SubtreeSize

	parentNewBups, err := q.InsertBottomUpProps(ctx, queries.InsertBottomUpPropsParams{
		Children:         parentbups.Children,
		SubtreeSize:      parentbups.SubtreeSize,
		SubtreeCompleted: parentbups.SubtreeCompleted,
	})
	if err != nil {
		return fmt.Errorf("InsertBottomUpProps: %w", err)
	}

	parentNew, err := q.InsertUpdatedObjective(ctx, queries.InsertUpdatedObjectiveParams{
		Oid:       parent.Oid,
		Based:     parent.Vid,
		CreatedBy: op.Opid,
		Pid:       parent.Pid,
		Bupid:     parentNewBups.Bupid,
	})
	if err != nil {
		return fmt.Errorf("InsertUpdatedObjective: %w", err)
	}

	sublinks, err := q.SelectSubLinks(ctx, queries.SelectSubLinksParams{
		SupOid: parent.Oid,
		SupVid: parent.Vid,
	})
	if err != nil {
		return fmt.Errorf("SelectSubLinks: %w", err)
	}
	for _, sublink := range sublinks {
		if sublink.SubOid == params.Subject.Oid { // deleted
			continue
		}
		_, err := q.InsertUpdatedLink(ctx, queries.InsertUpdatedLinkParams{
			SupOid:            parentNew.Oid,
			SupVid:            parentNew.Vid,
			SubOid:            sublink.SubOid,
			SubVid:            sublink.SubVid,
			CreatedAtOriginal: sublink.CreatedAtOriginal,
		})
		if err != nil {
			return fmt.Errorf("InsertUpdatedLink: %w", err)
		}
	}

	_, err = q.UpdateActiveVidForObjective(ctx, queries.UpdateActiveVidForObjectiveParams{
		Oid: parent.Oid,
		Vid: parentNew.Vid,
	})
	if err != nil {
		return fmt.Errorf("InsertActiveVidForObjective: %w", err)
	}

	ap[1] = models.Ovid{Oid: parentNew.Oid, Vid: parentNew.Vid}
	_, err = a.bubblink(ctx, q, ap[1:], op, deltas)
	if err != nil {
		return fmt.Errorf("bubblink: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
