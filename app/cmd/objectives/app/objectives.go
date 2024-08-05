package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateSubtaskParams struct {
	Creator columns.UserId
	Parent  models.Ovid
	Content string
}

// TODO: check prileges on parent
// DONE: create operations
// TODO: trigger task-props calculation
// TODO: transaction-commit-rollback
// DONE: bubblink
// DONE: mark active version for promoted ascendants
func (a *App) CreateSubtask(ctx context.Context, params CreateSubtaskParams) error {
	activepath, err := a.ListActivePathToRock(ctx, params.Parent)
	if err == ErrLeftBehind {
		return ErrLeftBehind
	} else if err != nil {
		return fmt.Errorf("checking the parent %s if it is inside active path: %w", params.Parent, err)
	}

	op, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.Parent.Oid,
		Subjectvid: params.Parent.Vid,
		Actor:      params.Creator,
		OpType:     database.OpTypeObjCreateSubtask,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("inserting the creation operation: %w", err)
	}

	_, err = a.queries.InsertOpObjCreateSubtask(ctx, database.InsertOpObjCreateSubtaskParams{
		Opid:    op.Opid,
		Content: pgtype.Text{String: params.Content, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("inserting subtask creation details: %w", err)
	}

	obj, err := a.queries.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: op.Opid,
		Props:     nil,
	})
	if err != nil {
		return fmt.Errorf("inserting the objective: %w", err)
	}

	// TODO: trigger computing props (async?)

	// bubblink: promote updates to ascendants
	child := models.Ovid{Oid: obj.Oid, Vid: obj.Vid}
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
