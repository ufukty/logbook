package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
)

// proposals can have multiple actions
// func (a *App) UpdateObjective(ovid Ovid, as []database.Action) error {
// 	nextVid := uuid.New()

// 	o, err := a.db.SelectObjective(ovid.Oid, ovid.Vid)
// 	if err != nil {
// 		return fmt.Errorf("getting objective for oid: %w", err)
// 	}

// 	if err := a.ApplyActionsOnVersionedObjective(o.Clone(), as); err != nil {
// 		return fmt.Errorf("applying action list on objective: %w", err)
// 	}

// 	createNextVersionOfParent := func(oid database.ObjectiveId) database.Objective {
// 		links , err := a.db.SelectTheUpperLink(Ovid{oid, vid})
// 		for _, link := range links {
// 			if link.Type == nil {
// 			}
// 		}
// 		// TODO: same version of update sibling
// 	}
// 	updateChildren := func(ovid Ovid) {

// 	}

// 	bq := endpoints.TagAssignRequest{}
// 	bs, err := bq.Send()
// 	if err != nil {
// 		return fmt.Errorf("copying tag records to new version of the task: %w", err)
// 	}

// 	return nil
// }

// func (a *App) ComposeView(root database.ObjectiveId, vid database.VersionId) (any, error) {

// }

// TODO: Turn the parent objective into a goal if it is currently a task
// TODO: create NewOperation
// TODO: trigger task-props calculation
// TODO: transaction-commit-rollback
func (a *App) createVersionedObjective(ctx context.Context, act CreateObjectiveAction, ancestry []Ovid, vancestry []database.VersioningConfig) ([]Ovid, error) {
	// check authz
	vc := vancestry[len(vancestry)-1]
	v, err := a.queries.InsertVersion(ctx, vc.Effective)
	if err != nil {
		return nil, fmt.Errorf("producing the next version id before updating ancestry: %w", err)
	}

	o := database.Objective{
		Oid:     act.Parent.Oid,
		Vid:     v.Vid,
		Based:   database.ZeroVersionId,
		Content: act.Content,
		Creator: act.Creator,
	}
	o, err = a.queries.InsertObjective(ctx, database.InsertObjectiveParams{
		Vid:     o.Vid,
		Based:   o.Based,
		Content: o.Content,
		Creator: o.Creator,
	})
	if err != nil {
		return nil, fmt.Errorf("inserting objective into the db: %w", err)
	}

	updates := []Ovid{}
	var prev Ovid
	for _, parentOvid := range ancestry {
		parent, err := a.queries.SelectObjective(ctx, database.SelectObjectiveParams{
			Oid: parentOvid.Oid,
			Vid: parentOvid.Vid,
		})
		if err != nil {
			return nil, fmt.Errorf("selecting parent %s from db: %w", parentOvid, err)
		}
		parent.Based, parent.Vid = parent.Vid, v.Vid
		parent, err = a.queries.InsertObjective(ctx, database.InsertObjectiveParams{
			Vid:     parent.Vid,
			Based:   parent.Based,
			Content: parent.Content,
			Creator: parent.Creator,
		})
		if err != nil {
			return nil, fmt.Errorf("inserting version bumped parent into the db: %w", err)
		}
		sublinks, err := a.queries.SelectSubLinks(ctx, database.SelectSubLinksParams{
			SupOid: parent.Oid,
			SupVid: parent.Vid,
		})
		if err != nil {
			return nil, fmt.Errorf("selecting sublinks of parent (%q/%q) from db: %w", parent.Oid, parent.Vid, err)
		}
		for _, link := range sublinks {
			if link.SubOid == prev.Oid {
				_, err = a.queries.InsertLink(ctx, database.InsertLinkParams{
					SupOid: parent.Oid,
					SupVid: parent.Vid,
					SubOid: prev.Oid,
					SubVid: v.Vid,
				})
				if err != nil {
					return nil, fmt.Errorf("inserting the link from %q (direct ancestry) to %q in version %q: %w", prev.Oid, parent.Oid, v.Vid, err)
				}
			} else {
				_, err = a.queries.InsertLink(ctx, database.InsertLinkParams{
					SupOid: parent.Oid,
					SupVid: parent.Vid,
					SubOid: link.SubOid,
					SubVid: link.SubVid,
				})
				if err != nil {
					return nil, fmt.Errorf("inserting the link from %q (sibling from ancestry) to %q in version %q: %w", prev.Oid, parent.Oid, v.Vid, err)
				}
			}
		}
		prev.Oid, prev.Vid = parent.Oid, parent.Vid
	}

	return updates, nil
}

func (a *App) CreateObjective(ctx context.Context, act CreateObjectiveAction) ([]Ovid, error) {
	ancestry, err := a.ListObjectiveAncestry(ctx, act.Parent)
	if err != nil {
		return nil, fmt.Errorf("listing ancestry of %q: %w", act.Parent, err)
	}
	vancestry, err := a.ListVersioningConfigForAncestry(ctx, ancestry)
	if err != nil {
		return nil, fmt.Errorf("listing versioning config for parents: %w", err)
	}
	if len(vancestry) > 0 {
		ovids, err := a.createVersionedObjective(ctx, act, ancestry, vancestry)
		if err != nil {
			return nil, fmt.Errorf("creating objective under versioning: %w", err)
		}
		return ovids, nil
	} else {
		return nil, nil
	}
}
