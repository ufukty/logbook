package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
)

// promotes an update to ascendants
// first item of the activepath should be the source of update promotion, newly updated objective (oid:latestvid)
// it returns the operation id generated for the transitive update of the uppermost objective in the activepath
func (a *App) bubblink(ctx context.Context, activepath []models.Ovid, op database.Operation) (columns.OperationId, error) {
	if len(activepath) <= 1 {
		return columns.ZeroOperationId, nil // no ascendant to promote update
	}
	child := activepath[0]
	cause := op.Opid
	for _, ascendant := range activepath[1:] {
		optrs, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
			Subjectoid: ascendant.Oid,
			Subjectvid: ascendant.Vid,
			Actor:      columns.ZeroUserId, // inherit user?
			OpType:     database.OpTypeTransitive,
			OpStatus:   database.OpStatusAccepted,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("inserting operation on ascendant %s for transitive update: %w", ascendant, err)
		}

		_, err = a.queries.InsertOpTransitive(ctx, database.InsertOpTransitiveParams{
			Opid:  optrs.Opid,
			Cause: cause,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("inserting transitive update specific operation details on ascendant %s for transitive update: %w", ascendant, err)
		}

		obj, err := a.queries.SelectObjective(ctx, database.SelectObjectiveParams{
			Oid: ascendant.Oid,
			Vid: ascendant.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("selecting current version of objective for props: %w", err)
		}

		objasc, err := a.queries.InsertUpdatedObjective(ctx, database.InsertUpdatedObjectiveParams{
			Oid:       ascendant.Oid,
			Based:     ascendant.Vid,
			CreatedBy: cause,
			Props:     obj.Props,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("inserting version updated ascendant: %w", err)
		}

		_, err = a.queries.InsertLink(ctx, database.InsertLinkParams{
			SupOid: objasc.Oid,
			SupVid: objasc.Vid,
			SubOid: child.Oid,
			SubVid: child.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("inserting a link the updated ascendants: %w", err)
		}

		// link unchanged siblings too
		sublinks, err := a.queries.SelectSubLinks(ctx, database.SelectSubLinksParams{
			SupOid: ascendant.Oid,
			SupVid: ascendant.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("selecting list of sub links of ascendant for current version: %w", err)
		}
		for _, sublink := range sublinks {
			if sublink.SubOid == objasc.Oid {
				continue // not sibling but itself's old version
			}
			_, err = a.queries.InsertLink(ctx, database.InsertLinkParams{
				SupOid: objasc.Oid,
				SupVid: objasc.Vid,
				SubOid: sublink.SubOid,
				SubVid: sublink.SubVid,
			})
			if err != nil {
				return columns.ZeroOperationId, fmt.Errorf("inserting a link from updated ascendants to existing sibling: %w", err)
			}
		}

		// TODO: trigger computing props on objasc (async?)

		// TODO: publish an event to notify frontends viewing any objective of active path?
		_, err = a.queries.UpdateActiveVidForObjective(ctx, database.UpdateActiveVidForObjectiveParams{
			Oid: objasc.Oid,
			Vid: objasc.Vid,
		})
		if err != nil {
			return columns.ZeroOperationId, fmt.Errorf("updating active version for the ascendant: %w", err)
		}
		cause = optrs.Opid
		child = models.Ovid{Oid: objasc.Oid, Vid: objasc.Vid}
	}
	return cause, nil
}
