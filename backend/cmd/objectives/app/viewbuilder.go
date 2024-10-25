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

var ViewportLimit = 200

var ErrGiantViewport = fmt.Errorf("requested portion of document is too big")

type ViewBuilderParams struct {
	Viewer        columns.UserId
	Root          models.Ovid
	Start, Length int
}

type line struct{ start, end int32 }

func doOverlap(a, b line) bool {
	return a.start < b.end && b.start < a.end
}

func (a *App) isFold(ctx context.Context, q *database.Queries, viewer columns.UserId, subject columns.ObjectiveId) (bool, error) {
	vp, err := q.SelectObjectiveViewPrefs(ctx, database.SelectObjectiveViewPrefsParams{
		Uid: viewer,
		Oid: subject,
	})
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("SelectObjectiveViewPrefs: %w", err)
	} else {
		return vp.Fold, nil
	}
}

// FIXME: apply recursion as permissions allow for solo and collaborated objectives
func (a *App) getUss(ctx context.Context, q *database.Queries, viewer columns.UserId, subject models.Ovid) (int32, error) {
	if cachedsize, ok := a.caches.Uss.Get(usssubject{Viewer: viewer, Object: subject}); ok {
		return cachedsize, nil
	}
	subs, err := q.SelectSubLinks(context.Background(), database.SelectSubLinksParams{
		SupOid: subject.Oid,
		SupVid: subject.Vid,
	})
	if err != nil {
		return 0, fmt.Errorf("SelectSubLinks: %w", err)
	}
	size := int32(0)
	for _, sub := range subs {
		fold, err := a.isFold(ctx, q, viewer, sub.SubOid)
		if err != nil {
			return 0, fmt.Errorf("isFold: %w", err)
		}
		if fold {
			size += 1
		} else {
			size_, err := a.getUss(ctx, q, viewer, models.Ovid{Oid: sub.SubOid, Vid: sub.SubVid})
			if err != nil {
				return 0, fmt.Errorf("getUss(%s): %w", sub.SubOid, err)
			}
			size += size_ + 1
		}
	}
	a.caches.Uss.Set(usssubject{Viewer: viewer, Object: subject}, size)
	return size, nil
}

// TODO: mind permissions
func (a *App) viewBuilder(ctx context.Context, q *database.Queries, viewer columns.UserId, subject models.Ovid, start, end int32, depth int) ([]owners.DocumentItem, error) {
	doc := []owners.DocumentItem{}

	// if end <= start || end < 0 {
	// 	return doc, nil
	// }

	unfoldSubtreeSize, err := a.getUss(ctx, q, viewer, subject)
	if err != nil {
		return nil, fmt.Errorf("getUss/1: %w", err)
	}

	// if unfoldSubtreeSize+1 < start { // viewport starts after the subtree
	// 	return doc, nil
	// }

	fold, err := a.isFold(ctx, q, viewer, subject.Oid)
	if err != nil {
		return nil, fmt.Errorf("isFold: %w", err)
	}

	if start <= 0 {
		doc = append(doc, owners.DocumentItem{
			Oid:           subject.Oid,
			Vid:           subject.Vid,
			Depth:         depth,
			ObjectiveType: ternary(unfoldSubtreeSize > 0, owners.Goal, owners.Task),
			Folded:        fold,
		})
	}

	if !fold && doOverlap(line{1, unfoldSubtreeSize + 1}, line{start, end}) {
		subs, err := q.SelectSubLinks(ctx, database.SelectSubLinksParams{
			SupOid: subject.Oid,
			SupVid: subject.Vid,
		})
		if err != nil {
			return nil, fmt.Errorf("SelectSubLinks: %w", err)
		}
		passed := int32(1)
		for _, sub := range subs {
			subUss, err := a.getUss(ctx, q, viewer, models.Ovid{Oid: sub.SubOid, Vid: sub.SubVid})
			if err != nil {
				return nil, fmt.Errorf("getUss: %w", err)
			}
			if doOverlap(line{passed, passed + subUss + 1}, line{start, end}) {
				v, err := a.viewBuilder(ctx, q, viewer, models.Ovid{sub.SubOid, sub.SubVid}, start-passed, end-passed, depth+1)
				if err != nil {
					return nil, fmt.Errorf("viewBuilder(%s): %w", sub.SubOid, err)
				}
				if len(v) > 0 {
					doc = slices.Concat(doc, v)
				}
			}
			passed += subUss + 1
		}
	}

	return doc, nil
}

// [App.ViewBuilder] will return an array of [models.Ovid] belong to objectives sit
// in the boundaries of viewport described as a range starts from [ViewBuilderParams.Start]
// and in the length of [ViewBuilderParams.Length] which doesn't contain hidden/unaccessible objectives.
//
// may wrap: [l2.ErrLeftBehind]
func (a *App) ViewBuilder(ctx context.Context, params ViewBuilderParams) ([]owners.DocumentItem, error) {
	if params.Length > ViewportLimit {
		return nil, ErrGiantViewport
	}

	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	v, err := a.viewBuilder(ctx, q, params.Viewer, params.Root, int32(params.Start), int32(params.Start+params.Length), 0)
	if err != nil {
		return nil, fmt.Errorf("viewBuilder: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return v, nil
}
