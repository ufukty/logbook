package l2

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/internal/logger"
	"logbook/models"
	"slices"
)

var ErrLeftBehind = fmt.Errorf("the objective is either directly or eventually linked to an objective which its version is left behind")

func helperListActivePathToRock(ctx context.Context, q *database.Queries, subject models.Ovid, l *logger.Logger) ([]models.Ovid, error) {
	active, err := q.SelectActive(ctx, subject.Oid)
	if err != nil {
		return nil, fmt.Errorf("SelectActive(%s): %w", subject.Oid, err)
	}
	if subject.Vid != active.Vid {
		return nil, ErrLeftBehind
	}

	activeparents, err := q.SelectUpperLinksToActiveObjectives(ctx, database.SelectUpperLinksToActiveObjectivesParams{
		SubOid: subject.Oid,
		SubVid: subject.Vid,
	})
	if err != nil {
		return nil, fmt.Errorf("SelectUpperLinksToActiveObjectives(%s): %w", subject, err)
	}
	if len(activeparents) == 0 {
		return []models.Ovid{subject}, nil // assuming it is the rock or an orphaned subtree
	}
	if len(activeparents) == 50 {
		l.Println("caution: helperListActivePathToRock has reached to the limit of 50 for selecting upper links")
	}

	for _, activeparent := range activeparents {
		activeparentovid := models.Ovid{
			Oid: activeparent.SupOid,
			Vid: activeparent.SupVid,
		}
		path, err := helperListActivePathToRock(ctx, q, activeparentovid, l)
		if err == ErrLeftBehind {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("helperListActivePathToRock(%s): %w", activeparentovid, err)
		} else {
			return append(path, subject), nil
		}
	}
	return nil, ErrLeftBehind
}

// returns a path from the node to the root: [subject, active-ascendants..., rock] or [ErrLeftBehind]
func (a *App) ListActivePathToRock(ctx context.Context, q *database.Queries, subject models.Ovid) ([]models.Ovid, error) {
	ap, err := helperListActivePathToRock(ctx, q, subject, a.l)
	if err != nil {
		return nil, fmt.Errorf("helperListActivePathToRock: %w", err)
	}
	slices.Reverse(ap)
	return ap, nil
}
