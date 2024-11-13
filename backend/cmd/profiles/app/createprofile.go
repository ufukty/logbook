package app

import (
	"context"
	"fmt"
	"logbook/cmd/profiles/database"
	"logbook/models/columns"
)

type CreateProfileParams struct {
	Uid       columns.UserId
	Firstname columns.HumanName
	Lastname  columns.HumanName
}

func (a *App) CreateProfile(ctx context.Context, params CreateProfileParams) error {
	_, err := a.oneshot.InsertProfileInformation(ctx, database.InsertProfileInformationParams{
		Uid:       params.Uid,
		Firstname: params.Firstname,
		Lastname:  params.Lastname,
	})
	if err != nil {
		return fmt.Errorf("insert profile information into db: %w", err)
	}
	return nil
}
