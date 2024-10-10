package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models/columns"
)

func (a *App) RockCreate(ctx context.Context, uid columns.UserId) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	props, err := q.InsertProperties(ctx, database.InsertPropertiesParams{
		Content:   "",
		Completed: false,
		Creator:   uid,
		Owner:     uid,
	})
	if err != nil {
		return fmt.Errorf("InsertProperties: %w", err)
	}

	bup, err := q.InsertBottomUpProps(ctx, database.InsertBottomUpPropsParams{
		Children:         0,
		SubtreeSize:      0,
		SubtreeCompleted: 0,
	})
	if err != nil {
		return fmt.Errorf("InsertBottomUpProps: %w", err)
	}

	obj, err := q.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: columns.ZeroOperationId,
		Pid:       props.Pid,
		Bupid:     bup.Bupid,
	})
	if err != nil {
		return fmt.Errorf("InsertNewObjective: %w", err)
	}

	_, err = q.InsertBookmark(ctx, database.InsertBookmarkParams{
		Uid:    uid,
		Oid:    obj.Oid,
		Title:  "",
		IsRock: true,
	})
	if err != nil {
		return fmt.Errorf("InsertBookmark: %w", err)
	}

	_, err = q.InsertActiveVidForObjective(ctx, database.InsertActiveVidForObjectiveParams{
		Oid: obj.Oid,
		Vid: obj.Vid,
	})
	if err != nil {
		return fmt.Errorf("UpdateActiveVidForObjective: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
