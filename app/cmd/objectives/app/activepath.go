package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/models"
	"slices"
)

var ErrLeftBehind = fmt.Errorf("the objective is either directly or eventually linked to an objective which its version is left behind")

func listActivePathToRockHelper(ctx context.Context, q *queries.Queries, subject models.Ovid) ([]models.Ovid, error) {
	active, err := q.SelectActive(ctx, subject.Oid)
	if err != nil {
		return nil, fmt.Errorf("SelectActive(%s): %w", subject.Oid, err)
	}
	if subject.Vid != active.Vid {
		return nil, ErrLeftBehind
	}

	parents, err := q.SelectUpperLinks(ctx, queries.SelectUpperLinksParams{
		SubOid: subject.Oid,
		SubVid: subject.Vid,
	})
	if len(parents) == 0 {
		return []models.Ovid{subject}, nil // assuming it is the rock or an orphaned subtree
	} else if err != nil {
		return nil, fmt.Errorf("SelectUpperLinks(%s): %w", subject, err)
	}

	for _, parent := range parents {
		path, err := listActivePathToRockHelper(ctx, q, models.Ovid{parent.SupOid, parent.SupVid})
		if err == ErrLeftBehind {
			continue // try other parents
		} else if err != nil {
			return nil, fmt.Errorf("listActivePathToRockHelper(%s): %w", parent, err)
		} else {
			return append(path, subject), nil
		}
	}
	return nil, ErrLeftBehind
}

// returns a path from the node to the root: [subject, active-ascendants..., rock] or [ErrLeftBehind]
func (a *App) listActivePathToRock(ctx context.Context, q *queries.Queries, subject models.Ovid) ([]models.Ovid, error) {
	ap, err := listActivePathToRockHelper(ctx, q, subject)
	if err != nil {
		return nil, fmt.Errorf("listActivePathToRockHelper: %w", err)
	}
	slices.Reverse(ap)
	return ap, nil
}
