package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5/pgtype"
)

func (a *App) RockCreate(ctx context.Context, uid columns.UserId) error {
	props, err := a.queries.InsertProperties(ctx, database.InsertPropertiesParams{
		Content: "",
		Creator: uid,
	})
	if err != nil {
		return fmt.Errorf("InsertProperties: %w", err)
	}

	obj, err := a.queries.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: columns.ZeroOperationId,
		Props:     props.Propid,
	})
	if err != nil {
		return fmt.Errorf("InsertNewObjective: %w", err)
	}

	_, err = a.queries.InsertBookmark(ctx, database.InsertBookmarkParams{
		Uid:         uid,
		Oid:         obj.Oid,
		Vid:         obj.Vid,
		DisplayName: pgtype.Text{"", true},
		IsRock:      true,
	})
	if err != nil {
		return fmt.Errorf("InsertBookmark: %w", err)
	}

	return nil
}
