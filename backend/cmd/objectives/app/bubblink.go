package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
)

// represents the amount of change in the Bottom Up Props of each ascendant
type bubblinkDeltaValues struct {
	SubtreeCompleted int32
	SubtreeSize      int32
}

var zeroDeltas bubblinkDeltaValues

// promotes an update to ascendants
// returns the uppermost objective's operation id
func (a *App) bubblink(ctx context.Context, q *database.Queries, activepath []models.Ovid, op database.Operation, deltas bubblinkDeltaValues) (columns.OperationId, error) {
	if len(activepath) <= 1 {
		return columns.ZeroOperationId, nil // no ascendant to promote update
	}
	child := activepath[0]
	cause := op.Opid
	for _, ascendant := range activepath[1:] {
		optrs, err := q.InsertOperation(ctx, database.InsertOperationParams{
			Subjectoid: ascendant.Oid,
			Subjectvid: ascendant.Vid,
			Actor:      columns.ZeroUserId, // inherit user?
			OpType:     database.OpTypeTransitive,
			OpStatus:   database.OpStatusAccepted,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("InsertOperation(%s): %w", ascendant, err)
		}

		_, err = q.InsertOpTransitive(ctx, database.InsertOpTransitiveParams{
			Opid:  optrs.Opid,
			Cause: cause,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("InsertOpTransitive(%s): %w", ascendant, err)
		}

		objAsc, err := q.SelectObjective(ctx, database.SelectObjectiveParams{
			Oid: ascendant.Oid,
			Vid: ascendant.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("SelectObjective: %w", err)
		}

		var bupid columns.BottomUpPropsId
		if deltas != zeroDeltas {
			bup, err := q.SelectBottomUpProps(ctx, objAsc.Bupid)
			if err != nil {
				return columns.ZeroOperationId, fmt.Errorf("SelectBottomUpProps: %w", err)
			}

			bup.SubtreeSize += deltas.SubtreeSize
			bup.SubtreeCompleted += deltas.SubtreeCompleted

			bupUpdated, err := q.InsertBottomUpProps(ctx, database.InsertBottomUpPropsParams{
				Children:         bup.Children,
				SubtreeSize:      bup.SubtreeSize,
				SubtreeCompleted: bup.SubtreeCompleted,
			})
			if err != nil {
				return columns.ZeroOperationId, fmt.Errorf("InsertBottomUpProps: %w", err)
			}
			bupid = bupUpdated.Bupid
		} else {
			bupid = objAsc.Bupid
		}

		objAscUpd, err := q.InsertUpdatedObjective(ctx, database.InsertUpdatedObjectiveParams{
			Oid:       ascendant.Oid,
			Based:     ascendant.Vid,
			CreatedBy: optrs.Opid,
			Pid:       objAsc.Pid,
			Bupid:     bupid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("InsertUpdatedObjective: %w", err)
		}

		// link unchanged siblings too
		sublinks, err := q.SelectSubLinks(ctx, database.SelectSubLinksParams{
			SupOid: ascendant.Oid,
			SupVid: ascendant.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("SelectSubLinks: %w", err)
		}
		for _, sublink := range sublinks {
			vid := sublink.SubVid
			if sublink.SubOid == child.Oid { // not a sibling, but updated child's old version
				vid = child.Vid
			}
			_, err = q.InsertUpdatedLink(ctx, database.InsertUpdatedLinkParams{
				SupOid:            objAscUpd.Oid,
				SupVid:            objAscUpd.Vid,
				SubOid:            sublink.SubOid,
				SubVid:            vid,
				CreatedAtOriginal: sublink.CreatedAtOriginal,
			})
			if err != nil {
				return columns.ZeroOperationId, fmt.Errorf("InsertUpdatedLink: %w", err)
			}
		}

		// TODO: publish an event to notify frontends viewing any objective of active path?
		_, err = q.UpdateActiveVidForObjective(ctx, database.UpdateActiveVidForObjectiveParams{
			Oid: objAscUpd.Oid,
			Vid: objAscUpd.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("UpdateActiveVidForObjective: %w", err)
		}
		cause = optrs.Opid
		child = models.Ovid{Oid: objAscUpd.Oid, Vid: objAscUpd.Vid}
	}

	return cause, nil
}
