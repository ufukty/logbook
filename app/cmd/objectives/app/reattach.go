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
	ComesAfter    columns.ObjectiveId
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

// TODO: check auth at the both current and next parent for actor
func (a *App) Reattach(ctx context.Context, params ReattachParams) error {
	apCurrent, err := a.listActivePathToRock(ctx, params.CurrentParent)
	if err != nil {
		return fmt.Errorf("checking if the current parent is in active path: %w", err)
	}

	apNext, err := a.listActivePathToRock(ctx, params.NextParent)
	if err != nil {
		return fmt.Errorf("checking if the next parent is in active path: %w", err)
	}

	opDetach, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.CurrentParent.Oid,
		Subjectvid: params.CurrentParent.Vid,
		Actor:      params.Actor,
		OpType:     database.OpTypeObjDetach,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("inserting operation for detaching the objective from current parrent: %w", err)
	}

	_, err = a.queries.InsertOpObjAttach(ctx, database.InsertOpObjAttachParams{
		Opid:  opDetach.Opid,
		Child: params.CurrentParent.Oid,
	})
	if err != nil {
		return fmt.Errorf("inserting operation specific details for detaching the objective from current parent: %w", err)
	}

	opAttach, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.NextParent.Oid,
		Subjectvid: params.NextParent.Vid,
		Actor:      opDetach.Actor,
		OpType:     database.OpTypeObjAttach,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("inserting operation for attaching the objective to next parrent: %w", err)
	}

	_, err = a.queries.InsertOpObjAttach(ctx, database.InsertOpObjAttachParams{
		Opid:  opAttach.Opid,
		Child: params.NextParent.Oid,
	})
	if err != nil {
		return fmt.Errorf("inserting operation specific details for attaching the objective to next parent: %w", err)
	}

	apCurrent, apNext, apCommon := popCommonActivePath(apCurrent, apNext)
	opidCurrent, err := a.bubblink(ctx, apCurrent, opDetach)
	if err != nil {
		return fmt.Errorf("promoting the update to the the current parent is in active path: %w", err)
	}
	opidNext, err := a.bubblink(ctx, apNext, opAttach)
	if err != nil {
		return fmt.Errorf("promoting the update to the the next parent is in active path: %w", err)
	}
	if len(apCommon) > 0 {
		opMerg, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
			Subjectoid: apCommon[len(apCommon)-1].Oid,
			Subjectvid: apCommon[len(apCommon)-1].Vid,
			Actor:      columns.ZeroUserId,
			OpType:     database.OpTypeDoubleTransitiveMerger,
			OpStatus:   database.OpStatusAccepted,
		})
		if err != nil {
			return fmt.Errorf("inserting operation for merging attachment and detachment operations that crossed on the same ascendant in their bubblinks: %w", err)
		}
		_, err = a.queries.InsertOpDoubleTransitiveMerger(ctx, database.InsertOpDoubleTransitiveMergerParams{
			Opid:   opMerg.Opid,
			First:  opidCurrent,
			Second: opidNext,
		})
		if err != nil {
			return fmt.Errorf("inserting operation specific details for merging attachment and detachment operations that crossed on the same ascendant in their bubblinks: %w", err)
		}
		_, err = a.bubblink(ctx, apCommon, opAttach)
		if err != nil {
			return fmt.Errorf("promoting the update to the overlapping active paths: %w", err)
		}
	}

	return nil
}
