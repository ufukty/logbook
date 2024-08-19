package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models/columns"
)

func (a *App) RockCreate(ctx context.Context, uid columns.UserId) error {
	props, err := a.queries.InsertProperties(ctx, database.InsertPropertiesParams{
		Content:   "",
		Completed: false,
		Creator:   uid,
		Owner:     uid,
	})
	if err != nil {
		return fmt.Errorf("InsertProperties: %w", err)
	}

	bup, err := a.queries.InsertBottomUpProps(ctx, database.InsertBottomUpPropsParams{
		SubtreeSize:      0,
		SubtreeCompleted: 0,
	})
	if err != nil {
		return fmt.Errorf("InsertBottomUpProps: %w", err)
	}

	obj, err := a.queries.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: columns.ZeroOperationId,
		Pid:       props.Pid,
		Bupid:     bup.Bupid,
	})
	if err != nil {
		return fmt.Errorf("InsertNewObjective: %w", err)
	}

	_, err = a.queries.InsertBookmark(ctx, database.InsertBookmarkParams{
		Uid:    uid,
		Oid:    obj.Oid,
		Title:  "",
		IsRock: true,
	})
	if err != nil {
		return fmt.Errorf("InsertBookmark: %w", err)
	}

	_, err = a.queries.InsertActiveVidForObjective(ctx, database.InsertActiveVidForObjectiveParams{
		Oid: obj.Oid,
		Vid: obj.Vid,
	})
	if err != nil {
		return fmt.Errorf("UpdateActiveVidForObjective: %w", err)
	}

	return nil
}
