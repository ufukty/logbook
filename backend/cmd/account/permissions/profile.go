package permissions

import (
	"context"
	"logbook/models/columns"
)

func (d Decider) CanUserSetProfile(ctx context.Context, user, profile columns.UserId) error {
	if user != profile {
		return ErrUnauthorized
	}
	return nil
}
