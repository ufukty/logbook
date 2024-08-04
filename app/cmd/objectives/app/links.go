package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
	"slices"

	"github.com/jackc/pgx/v5"
)

var ErrLeftBehind = fmt.Errorf("the objective is either directly or eventually linked to an objective which its version is left behind")

// returns a path from the node to the root: [subject, active-ascendants, ..., rock] or [ErrLeftBehind]
func (a *App) ListActivePathToRock(ctx context.Context, subject models.Ovid) ([]models.Ovid, error) {
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
	if err == pgx.ErrNoRows {
		return []models.Ovid{subject}, nil // assuming it is the rock or an orphaned subtree
	} else if err != nil {
		return nil, fmt.Errorf("SelectUpperLinks(%s): %w", subject, err)
	}

	for _, parent := range parents {
		path, err := a.ListActivePathToRock(ctx, models.Ovid{parent.SupOid, parent.SupVid})
		if err == ErrLeftBehind {
			continue // try other parents
		} else if err != nil {
			return nil, fmt.Errorf("ListActivePathToRock(%s): %w", parent, err)
		} else {
			return slices.Insert(path, 0, subject), nil
		}
	}
	return nil, ErrLeftBehind
}

func (a *App) ListObjectiveAncestry(ctx context.Context, ovid models.Ovid) ([]models.Ovid, error) {
	anc := []models.Ovid{}
	c := ovid
	for limit := 0; true; limit++ {
		l, err := a.queries.SelectTheUpperLink(ctx, database.SelectTheUpperLinkParams{
			SubOid: c.Oid,
			SubVid: c.Vid,
		})
		if err != nil {
			return nil, fmt.Errorf("db.SelectTheUpperLink(%q, %q): %w", c.Oid, c.Vid, err)
		}
		c.Oid = l.SupOid
		c.Vid = l.SupVid
		if c.Oid == columns.ZeroObjectId {
			break
		}
		anc = append(anc, c)
		if limit == 100 {
			return nil, fmt.Errorf("depth limit (100) is reached")
		}
	}
	return anc, nil
}
