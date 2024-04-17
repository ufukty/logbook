package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
)

type CreateProfileParams struct {
	SessionToken database.SessionToken
	Uid          database.UserId
	Firstname    database.HumanName
	Lastname     database.HumanName
}

func (a App) CreateProfile(ctx context.Context, params CreateProfileParams) error {
	err := a.authz.AssertCanSetProfile(ctx, params.SessionToken, params.Uid)
	if err != nil {
		return fmt.Errorf("checking authorization: %w", err)
	}

	_, err = a.queries.InsertProfileInformation(ctx, database.InsertProfileInformationParams{
		Uid:       params.Uid,
		Firstname: string(params.Firstname),
		Lastname:  string(params.Lastname),
	})
	if err != nil {
		return fmt.Errorf("insert profile information into db: %w", err)
	}
	return nil
}
