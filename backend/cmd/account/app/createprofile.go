package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
	"logbook/models/columns"
)

type CreateProfileParams struct {
	SessionToken columns.SessionToken
	Uid          columns.UserId
	Firstname    columns.HumanName
	Lastname     columns.HumanName
}

func (a App) CreateProfile(ctx context.Context, params CreateProfileParams) error {
	uid, err := a.s.UserForSessionToken(ctx, params.SessionToken)
	if err != nil {
		return fmt.Errorf("UserForSessionToken: %w", err)
	}

	err = a.pd.CanUserSetProfile(ctx, uid, params.Uid)
	if err != nil {
		return fmt.Errorf("checking authorization: %w", err)
	}

	_, err = a.oneshot.InsertProfileInformation(ctx, database.InsertProfileInformationParams{
		Uid:       params.Uid,
		Firstname: string(params.Firstname),
		Lastname:  string(params.Lastname),
	})
	if err != nil {
		return fmt.Errorf("insert profile information into db: %w", err)
	}
	return nil
}
