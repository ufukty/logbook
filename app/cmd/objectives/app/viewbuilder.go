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

// TODO: mind permissions
func (a *App) viewBuilder(ctx context.Context, viewer columns.UserId, subject models.Ovid, start, end int32, depth int) ([]owners.ObjectiveView, error) {
	view := []owners.ObjectiveView{}

	if end <= start || end < 0 {
		return view, nil
	}

	computedToTop, err := a.queries.SelectComputedToTop(ctx, database.SelectComputedToTopParams{
		Oid:    subject.Oid,
		Vid:    subject.Vid,
		Viewer: viewer,
	})
	if err != nil {
		return nil, fmt.Errorf("SelectComputedToTop/1: %w", err)
	}
	if computedToTop.SubtreeSize+1 < start { // viewport starts after the subtree
		return view, nil
	}

	subs, _ := a.queries.SelectSubLinks(ctx, database.SelectSubLinksParams{
		SupOid: subject.Oid,
		SupVid: subject.Vid,
	})

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
		view = append(view, owners.ObjectiveView{
			Oid:           subject.Oid,
			Vid:           subject.Vid,
			Depth:         depth,
			ObjectiveType: objtype,
			Folded:        fold,
		})
	}
	cursor++

	if fold {
		cursor += computedToTop.SubtreeSize
	} else {
		for _, sub := range subs {
			computedToTop, err := a.queries.SelectComputedToTop(ctx, database.SelectComputedToTopParams{
				Oid:    sub.SubOid,
				Vid:    sub.SubVid,
				Viewer: viewer,
			})
			if err != nil {
				return nil, fmt.Errorf("SelectComputedToTop/2: %w", err)
			}

			if doOverlap(start, end, cursor, cursor+computedToTop.SubtreeSize+1) {
				v, err := a.viewBuilder(ctx, viewer, models.Ovid{sub.SubOid, sub.SubVid}, start-cursor, end-cursor, depth+1)
				if err != nil {
					return nil, fmt.Errorf("viewBuilder: %w", err)
				}

				if len(v) > 0 {
					view = slices.Concat(view, v)
				}
			}
			cursor += computedToTop.SubtreeSize + 1
		}
	}

	return view, nil
}

// [App.ViewBuilder] will return an array of [models.Ovid] belong to objectives sit
// in the boundaries of viewport described as a range starts from [ViewBuilderParams.Start]
// and in the length of [ViewBuilderParams.Length] which doesn't contain hidden/unaccessible objectives.
//
// may wrap: [ErrLeftBehind]
func (a *App) ViewBuilder(ctx context.Context, params ViewBuilderParams) ([]owners.ObjectiveView, error) {
	return a.viewBuilder(ctx, params.Viewer, params.Root, int32(params.Start), int32(params.Start+params.Length), 0)
}
