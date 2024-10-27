package adapter

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5"
)

type Collaboration struct {
	q *database.Queries

	Coid columns.CollaborationId
}

// TODO: check access through group type collaborators
func (c *Collaboration) IsUserACollaborator(ctx context.Context, uid columns.UserId) (bool, error) {
	_, err := c.q.SelectUserTypeCollaboratorByUserId(ctx, database.SelectUserTypeCollaboratorByUserIdParams{
		Coid: c.Coid,
		Uid:  uid,
	})
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("SelectCollaborators: %w", err)
	}

	return true, nil
}
