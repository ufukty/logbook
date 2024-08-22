package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/models"
	"logbook/models/columns"
)

// represents the amount of change in the Bottom Up Props of each ascendant
type bubblinkDeltaValues struct {
	Children         int32 // only first parent
	SubtreeCompleted int32
	SubtreeSize      int32
}

var zeroDeltas bubblinkDeltaValues

// promotes an update to ascendants
// returns the uppermost objective's operation id
func (a *App) bubblink(ctx context.Context, q *queries.Queries, activepath []models.Ovid, op queries.Operation, deltas bubblinkDeltaValues) (columns.OperationId, error) {
	if len(activepath) <= 1 {
		return columns.ZeroOperationId, nil // no ascendant to promote update
	}
	child := activepath[0]
	cause := op.Opid
	for i, ascendant := range activepath[1:] {
		optrs, err := q.InsertOperation(ctx, queries.InsertOperationParams{
			Subjectoid: ascendant.Oid,
			Subjectvid: ascendant.Vid,
			Actor:      columns.ZeroUserId, // inherit user?
			OpType:     queries.OpTypeTransitive,
			OpStatus:   queries.OpStatusAccepted,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("InsertOperation(%s): %w", ascendant, err)
		}

		_, err = q.InsertOpTransitive(ctx, queries.InsertOpTransitiveParams{
			Opid:  optrs.Opid,
			Cause: cause,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("InsertOpTransitive(%s): %w", ascendant, err)
		}

		obj, err := q.SelectObjective(ctx, queries.SelectObjectiveParams{
			Oid: ascendant.Oid,
			Vid: ascendant.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("SelectObjective: %w", err)
		}

		var bupid columns.BottomUpPropsId
		if deltas != zeroDeltas {
			bup, err := q.SelectBottomUpProps(ctx, obj.Bupid)
			if err != nil {
				return columns.ZeroOperationId, fmt.Errorf("SelectBottomUpProps: %w", err)
			}

			bup.Children += ternary(i == 0, deltas.Children, 0)
			bup.SubtreeSize += deltas.SubtreeSize
			bup.SubtreeCompleted += deltas.SubtreeCompleted

			bupUpdated, err := q.InsertBottomUpProps(ctx, queries.InsertBottomUpPropsParams{
				Children:         bup.Children,
				SubtreeSize:      bup.SubtreeSize,
				SubtreeCompleted: bup.SubtreeCompleted,
			})
			if err != nil {
				return columns.ZeroOperationId, fmt.Errorf("InsertBottomUpProps: %w", err)
			}
			bupid = bupUpdated.Bupid
		} else {
			bupid = obj.Bupid
		}

		objasc, err := q.InsertUpdatedObjective(ctx, queries.InsertUpdatedObjectiveParams{
			Oid:       ascendant.Oid,
			Based:     ascendant.Vid,
			CreatedBy: cause,
			Pid:       obj.Pid,
			Bupid:     bupid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("InsertUpdatedObjective: %w", err)
		}

		_, err = q.InsertLink(ctx, queries.InsertLinkParams{
			SupOid: objasc.Oid,
			SupVid: objasc.Vid,
			SubOid: child.Oid,
			SubVid: child.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("InsertLink/1: %w", err)
		}

		// link unchanged siblings too
		sublinks, err := q.SelectSubLinks(ctx, queries.SelectSubLinksParams{
			SupOid: ascendant.Oid,
			SupVid: ascendant.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("SelectSubLinks: %w", err)
		}
		for _, sublink := range sublinks {
			if sublink.SubOid == objasc.Oid {
				continue // not sibling but itself's old version
			}
			_, err = q.InsertLink(ctx, queries.InsertLinkParams{
				SupOid: objasc.Oid,
				SupVid: objasc.Vid,
				SubOid: sublink.SubOid,
				SubVid: sublink.SubVid,
			})
			if err != nil {
				return columns.ZeroOperationId, fmt.Errorf("InsertLink/2: %w", err)
			}
		}

		// TODO: publish an event to notify frontends viewing any objective of active path?
		_, err = q.UpdateActiveVidForObjective(ctx, queries.UpdateActiveVidForObjectiveParams{
			Oid: objasc.Oid,
			Vid: objasc.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("UpdateActiveVidForObjective: %w", err)
		}
		cause = optrs.Opid
		child = models.Ovid{Oid: objasc.Oid, Vid: objasc.Vid}
	}

	return cause, nil
}
