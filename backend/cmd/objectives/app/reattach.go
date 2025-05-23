package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
)

type ReattachParams struct {
	Actor         columns.UserId
	CurrentParent models.Ovid
	NextParent    models.Ovid
	ComesAfter    columns.ObjectiveId // TODO: remove
	Subject       columns.ObjectiveId
}

func popCommonActivePath(l, r []models.Ovid) ([]models.Ovid, []models.Ovid, []models.Ovid) {
	if len(l) == 0 || len(r) == 0 {
		return l, r, []models.Ovid{}
	}
	if len(l) > len(r) {
		l, r = r, l
	}
	lc, rc := len(l)-1, len(r)-1
	for rc > 0 && l[lc] == r[rc] {
		lc--
		rc--
	}
	common := l[lc+1:]
	return l[:lc+1], r[:rc+1], common
}

func (a *App) deltaValuesForReattachment(ctx context.Context, q *database.Queries, obj database.Objective) (bubblinkDeltaValues, bubblinkDeltaValues, error) {
	deltasCurrent := bubblinkDeltaValues{}
	deltasNext := bubblinkDeltaValues{}
	props, err := q.SelectProperties(ctx, obj.Pid)
	if err != nil {
		return bubblinkDeltaValues{}, bubblinkDeltaValues{}, fmt.Errorf("SelectProperties: %w", err)
	}
	if props.Completed {
		deltasCurrent.SubtreeCompleted--
		deltasNext.SubtreeCompleted++
	} else {
		deltasCurrent.SubtreeCompleted++
		deltasNext.SubtreeCompleted--
	}
	bups, err := q.SelectBottomUpProps(ctx, obj.Bupid)
	if err != nil {
		return bubblinkDeltaValues{}, bubblinkDeltaValues{}, fmt.Errorf("SelectBottomUpProps: %w", err)
	}
	deltasCurrent.SubtreeCompleted -= bups.SubtreeCompleted
	deltasNext.SubtreeCompleted += bups.SubtreeCompleted
	deltasCurrent.SubtreeSize -= bups.SubtreeSize
	deltasNext.SubtreeSize += bups.SubtreeSize
	return deltasCurrent, deltasNext, nil
}

// TODO: check auth at the both current and next parent for actor
func (a *App) Reattach(ctx context.Context, params ReattachParams) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	apCurrent, err := a.l2.ListActivePathToRock(ctx, q, params.CurrentParent)
	if err != nil {
		return fmt.Errorf("listActivePathToRock/current: %w", err)
	}

	apNext, err := a.l2.ListActivePathToRock(ctx, q, params.NextParent)
	if err != nil {
		return fmt.Errorf("listActivePathToRock/next: %w", err)
	}

	opDetach, err := q.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.CurrentParent.Oid,
		Subjectvid: params.CurrentParent.Vid,
		Actor:      params.Actor,
		OpType:     database.OpTypeObjDetach,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("InsertOperation: %w", err)
	}

	_, err = q.InsertOpObjDetach(ctx, database.InsertOpObjDetachParams{
		Opid:  opDetach.Opid,
		Child: params.CurrentParent.Oid,
	})
	if err != nil {
		return fmt.Errorf("InsertOpObjDetach: %w", err)
	}

	opAttach, err := q.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.NextParent.Oid,
		Subjectvid: params.NextParent.Vid,
		Actor:      opDetach.Actor,
		OpType:     database.OpTypeObjAttach,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("InsertOperation: %w", err)
	}

	_, err = q.InsertOpObjAttach(ctx, database.InsertOpObjAttachParams{
		Opid:  opAttach.Opid,
		Child: params.NextParent.Oid,
	})
	if err != nil {
		return fmt.Errorf("InsertOpObjAttach: %w", err)
	}

	active, err := q.SelectActive(ctx, params.Subject)
	if err != nil {
		return fmt.Errorf("SelectActive: %w", err)
	}

	obj, err := q.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: params.Subject,
		Vid: active.Vid,
	})
	if err != nil {
		return fmt.Errorf("SelectObjective: %w", err)
	}

	deltasCurrent, deltasNext, err := a.deltaValuesForReattachment(ctx, q, obj)
	if err != nil {
		return fmt.Errorf("deltaValuesForReattachment: %w", err)
	}

	apCurrent, apNext, apCommon := popCommonActivePath(apCurrent, apNext)
	opidCurrent, err := a.bubblink(ctx, q, apCurrent, opDetach, deltasCurrent)
	if err != nil {
		return fmt.Errorf("bubblink/current: %w", err)
	}
	opidNext, err := a.bubblink(ctx, q, apNext, opAttach, deltasNext)
	if err != nil {
		return fmt.Errorf("bubblink/next: %w", err)
	}
	if len(apCommon) > 0 {
		opMerg, err := q.InsertOperation(ctx, database.InsertOperationParams{
			Subjectoid: apCommon[len(apCommon)-1].Oid,
			Subjectvid: apCommon[len(apCommon)-1].Vid,
			Actor:      columns.ZeroUserId,
			OpType:     database.OpTypeDoubleTransitiveMerger,
			OpStatus:   database.OpStatusAccepted,
		})
		if err != nil {
			return fmt.Errorf("InsertOperation/merge: %w", err)
		}
		_, err = q.InsertOpDoubleTransitiveMerger(ctx, database.InsertOpDoubleTransitiveMergerParams{
			Opid:   opMerg.Opid,
			First:  opidCurrent,
			Second: opidNext,
		})
		if err != nil {
			return fmt.Errorf("InsertOpDoubleTransitiveMerger: %w", err)
		}
		_, err = a.bubblink(ctx, q, apCommon, opAttach, zeroDeltas)
		if err != nil {
			return fmt.Errorf("bubblink: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
