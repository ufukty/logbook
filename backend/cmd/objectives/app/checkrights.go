package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5"
)

type canSeeParams struct {
	Viewer  columns.UserId
	Subject models.Ovid
}

type Right string

const (
	RightNone      = Right("none")
	RightRead      = Right("read")
	RightReadWrite = Right("read-write")
)

func (a *App) checkRights(ctx context.Context, params canSeeParams) (Right, error) {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return RightNone, fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	ca, err := a.findControlArea(ctx, params.Subject)
	if err != nil {
		return RightNone, fmt.Errorf("findControlArea: %w", err)
	}

	switch ca.ArType {
	case database.AreaTypeSolo:
		obj, err := q.SelectObjective(ctx, database.SelectObjectiveParams{
			Oid: params.Subject.Oid,
			Vid: params.Subject.Vid,
		})
		if err != nil {
			return RightNone, fmt.Errorf("SelectObjective: %w", err)
		}
		props, err := q.SelectProperties(ctx, obj.Pid)
		if err != nil {
			return RightNone, fmt.Errorf("SelectProperties: %w", err)
		}
		if props.Creator == params.Viewer {
			return RightReadWrite, nil
		} else {
			return RightNone, nil
		}

	case database.AreaTypeCollaboration:
		co, err := q.SelectCollaborationOnControlArea(ctx, ca.Aid)
		if err != nil {
			return RightNone, fmt.Errorf("SelectCollaboration: %w", err)
		}

		_, err = q.SelectCollaborator(ctx, database.SelectCollaboratorParams{
			Cid: co.Cid,
			Uid: params.Viewer,
		})
		if err == pgx.ErrNoRows {
			return RightNone, nil
		} else if err != nil {
			return RightNone, fmt.Errorf("SelectCollaborators: %w", err)
		} else {
			return RightReadWrite, nil
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return RightNone, fmt.Errorf("commit: %w", err)
	}

	return RightNone, nil
}
