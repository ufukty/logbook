package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
	"logbook/models/owners"
	"slices"

	"github.com/jackc/pgx/v5"
)

type ViewBuilderParams struct {
	Viewer        columns.UserId
	Root          models.Ovid
	Start, Length int
}

func doOverlap(aStart, aEnd, bStart, bEnd int32) bool {
	return aStart < bEnd && bStart < aEnd
}

func (a *App) getSubtreeSize(ctx context.Context, viewer columns.UserId, bupid columns.BottomUpPropsId) (int32, error) {
	bup, err := a.queries.SelectBottomUpProps(ctx, bupid)
	if err != nil {
		return 0, fmt.Errorf("SelectBottomUpProps: %w", err)
	}
	buptp, err := a.queries.SelectBottomUpPropsThirdPerson(ctx, database.SelectBottomUpPropsThirdPersonParams{
		Bupid:  bupid,
		Viewer: viewer,
	})
	if err != nil && err != pgx.ErrNoRows {
		return 0, fmt.Errorf("SelectBottomUpPropsThirdPerson: %w", err)
	} else if err == pgx.ErrNoRows {
		return bup.SubtreeSize, nil
	} else {
		return bup.SubtreeSize + buptp.SubtreeSize, nil
	}
}

// TODO: mind permissions
func (a *App) viewBuilder(ctx context.Context, viewer columns.UserId, subject models.Ovid, start, end int32, depth int) ([]owners.DocumentItem, error) {
	view := []owners.DocumentItem{}

	if end <= start || end < 0 {
		return view, nil
	}

	obj, err := a.queries.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: subject.Oid,
		Vid: subject.Vid,
	})
	if err != nil {
		return nil, fmt.Errorf("SelectObjective: %w", err)
	}

	subtreeSize, err := a.getSubtreeSize(ctx, viewer, obj.Bupid)
	if err != nil {
		return nil, fmt.Errorf("getSubtreeSize/1: %w", err)
	}

	if subtreeSize+1 < start { // viewport starts after the subtree
		return view, nil
	}

	subs, err := a.queries.SelectSubLinks(ctx, database.SelectSubLinksParams{
		SupOid: subject.Oid,
		SupVid: subject.Vid,
	})
	if err != nil {
		return nil, fmt.Errorf("SelectSubLinks: %w", err)
	}

	fold := false
	cursor := int32(0)
	if doOverlap(start, end, 0, 1) {
		objtype := owners.Task
		if len(subs) > 0 {
			objtype = owners.Goal
		}
		vp, err := a.queries.SelectObjectiveViewPrefs(ctx, database.SelectObjectiveViewPrefsParams{
			Uid: viewer,
			Oid: subject.Oid,
		})
		if err != nil && err != pgx.ErrNoRows {
			return nil, fmt.Errorf("SelectObjectiveViewPrefs: %w", err)
		} else if err == nil {
			fold = vp.Fold
		}
		view = append(view, owners.DocumentItem{
			Oid:           subject.Oid,
			Vid:           subject.Vid,
			Depth:         depth,
			ObjectiveType: objtype,
			Folded:        fold,
		})
	}
	cursor++

	if fold {
		cursor += subtreeSize
	} else {
		for _, sub := range subs {
			subobj, err := a.queries.SelectObjective(ctx, database.SelectObjectiveParams{
				Oid: sub.SubOid,
				Vid: sub.SubVid,
			})
			if err != nil {
				return nil, fmt.Errorf("SelectObjective: %w", err)
			}
			subtreeSize, err := a.getSubtreeSize(ctx, viewer, subobj.Bupid)
			if err != nil {
				return nil, fmt.Errorf("subtreeSize/2: %w", err)
			}

			if doOverlap(start, end, cursor, cursor+subtreeSize+1) {
				v, err := a.viewBuilder(ctx, viewer, models.Ovid{sub.SubOid, sub.SubVid}, start-cursor, end-cursor, depth+1)
				if err != nil {
					return nil, fmt.Errorf("viewBuilder: %w", err)
				}

				if len(v) > 0 {
					view = slices.Concat(view, v)
				}
			}
			cursor += subtreeSize + 1
		}
	}

	return view, nil
}

// [App.ViewBuilder] will return an array of [models.Ovid] belong to objectives sit
// in the boundaries of viewport described as a range starts from [ViewBuilderParams.Start]
// and in the length of [ViewBuilderParams.Length] which doesn't contain hidden/unaccessible objectives.
//
// may wrap: [ErrLeftBehind]
func (a *App) ViewBuilder(ctx context.Context, params ViewBuilderParams) ([]owners.DocumentItem, error) {
	return a.viewBuilder(ctx, params.Viewer, params.Root, int32(params.Start), int32(params.Start+params.Length), 0)
}
