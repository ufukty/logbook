package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
)

// promotes an update to ascendants
func (a *App) bubblink(ctx context.Context, subject database.Objective, op database.Operation, activepath []models.Ovid) error {
	child := models.Ovid{Oid: subject.Oid, Vid: subject.Vid}
	cause := op.Opid
	for i := len(activepath) - 1; i >= 0; i-- {
		ascendant := activepath[i]

		optrs, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
			Subjectoid: ascendant.Oid,
			Subjectvid: ascendant.Vid,
			Actor:      columns.ZeroUserId, // inherit user?
			OpType:     database.OpTypeTransitive,
			OpStatus:   database.OpStatusAccepted,
		})
		if err != nil {
			return fmt.Errorf("inserting operation on ascendant %s for transitive update: %w", ascendant, err)
		}

		_, err = a.queries.InsertOpTransitive(ctx, database.InsertOpTransitiveParams{
			Opid:  optrs.Opid,
			Cause: cause,
		})
		if err != nil {
			return fmt.Errorf("inserting transitive update specific operation details on ascendant %s for transitive update: %w", ascendant, err)
		}

		objasc, err := a.queries.InsertUpdatedObjective(ctx, database.InsertUpdatedObjectiveParams{
			Oid:       ascendant.Oid,
			Based:     ascendant.Vid,
			CreatedBy: cause,
			Props:     nil,
		})
		if err != nil {
			return fmt.Errorf("inserting version updated ascendant: %w", err)
		}

		_, err = a.queries.InsertLink(ctx, database.InsertLinkParams{
			SupOid: objasc.Oid,
			SupVid: objasc.Vid,
			SubOid: child.Oid,
			SubVid: child.Vid,
		})
		if err != nil {
			return fmt.Errorf("inserting a link the updated ascendants: %w", err)
		}

		// link unchanged siblings too
		sublinks, err := a.queries.SelectSubLinks(ctx, database.SelectSubLinksParams{
			SupOid: ascendant.Oid,
			SupVid: ascendant.Vid,
		})
		if err != nil {
			return fmt.Errorf("selecting list of sub links of ascendant for current version: %w", err)
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
				return fmt.Errorf("inserting a link from updated ascendants to existing sibling: %w", err)
			}
		}

		// TODO: trigger computing props on objasc (async?)

		// TODO: publish an event to notify frontends viewing any objective of active path?
		_, err = a.queries.UpdateActiveVidForObjective(ctx, database.UpdateActiveVidForObjectiveParams{
			Oid: objasc.Oid,
			Vid: objasc.Vid,
		})
		if err != nil {
			return fmt.Errorf("updating active version for the ascendant: %w", err)
		}
		cause = optrs.Opid
		child = models.Ovid{Oid: objasc.Oid, Vid: objasc.Vid}
	}
	return nil
}
