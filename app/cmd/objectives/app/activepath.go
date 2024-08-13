package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"slices"
)

var ErrLeftBehind = fmt.Errorf("the objective is either directly or eventually linked to an objective which its version is left behind")

// returns a path from the node to the root: [subject, active-ascendants..., rock] or [ErrLeftBehind]
func (a *App) listActivePathToRock(ctx context.Context, subject models.Ovid) ([]models.Ovid, error) {
	var helper func(ctx context.Context, subject models.Ovid) ([]models.Ovid, error)
	helper = func(ctx context.Context, subject models.Ovid) ([]models.Ovid, error) {
		active, err := a.queries.SelectActive(ctx, subject.Oid)
		if err != nil {
			return nil, fmt.Errorf("SelectActive(%s): %w", subject.Oid, err)
		}
		if subject.Vid != active.Vid {
			return nil, ErrLeftBehind
		}

		parents, err := a.queries.SelectUpperLinks(ctx, database.SelectUpperLinksParams{
			SubOid: subject.Oid,
			SubVid: subject.Vid,
		})
		if len(parents) == 0 {
			return []models.Ovid{subject}, nil // assuming it is the rock or an orphaned subtree
		} else if err != nil {
			return nil, fmt.Errorf("SelectUpperLinks(%s): %w", subject, err)
		}

		for _, parent := range parents {
			path, err := helper(ctx, models.Ovid{parent.SupOid, parent.SupVid})
			if err == ErrLeftBehind {
				continue // try other parents
			} else if err != nil {
				return nil, fmt.Errorf("ListActivePathToRock(%s): %w", parent, err)
			} else {
				return append(path, subject), nil
			}
		}
		return nil, ErrLeftBehind
	}
	ap, err := helper(ctx, subject)
	if err != nil {
		return nil, err
	}
	slices.Reverse(ap)
	return ap, nil
}
