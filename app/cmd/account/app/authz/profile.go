package authz

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
)

func (a Authorization) AssertCanSetProfile(ctx context.Context, token database.SessionToken, uid database.UserId) error {
	sid, err := a.queries.SelectSessionByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("%w: selecting session id for session token from database: %w", NoAuthorization, err)
	}
	if uid != sid.Uid {
		return UnderAuthorized
	}
	return nil
}
