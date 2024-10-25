package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/app/l2"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
	"slices"

	"github.com/jackc/pgx/v5"
)

type CheckoutParams struct {
	User    columns.UserId
	Subject models.Ovid
	To      columns.VersionId
}

var ErrVersionDoesNotExist = fmt.Errorf("given version of the objective doesn't exist")

func (a *App) calculateDeltasForTwoVersions(ctx context.Context, q *database.Queries, src, dst columns.BottomUpPropsId) (bubblinkDeltaValues, error) {
	srcbup, err := q.SelectBottomUpProps(ctx, src)
	if err != nil {
		return zeroDeltas, fmt.Errorf("SelectBottomUpProps/src: %w", err)
	}

	dstbup, err := q.SelectBottomUpProps(ctx, dst)
	if err != nil {
		return zeroDeltas, fmt.Errorf("SelectBottomUpProps/dst: %w", err)
	}

	return bubblinkDeltaValues{
		SubtreeCompleted: dstbup.SubtreeCompleted - srcbup.SubtreeCompleted,
		SubtreeSize:      dstbup.SubtreeSize - srcbup.SubtreeSize,
	}, nil
}

func (a *App) Checkout(ctx context.Context, params CheckoutParams) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	activepath, err := a.l2.ListActivePathToRock(ctx, q, params.Subject)
	if err == l2.ErrLeftBehind {
		return l2.ErrLeftBehind
	} else if err != nil {
		return fmt.Errorf("listActivePathToRock: %w", err)
	}

	op, err := q.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.Subject.Oid,
		Subjectvid: params.Subject.Vid,
		Actor:      params.User,
		OpType:     database.OpTypeCheckout,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("InsertOperation: %w", err)
	}

	_, err = q.InsertOpCheckout(ctx, database.InsertOpCheckoutParams{
		Opid: op.Opid,
		To:   params.To,
	})
	if err != nil {
		return fmt.Errorf("InsertOpCheckout: %w", err)
	}

	dstobj, err := q.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: op.Subjectoid,
		Vid: params.To,
	})
	if err == pgx.ErrNoRows {
		return ErrVersionDoesNotExist
	} else if err != nil {
		return fmt.Errorf("SelectObjective/dst: %w", err)
	}

	srcobj, err := q.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: op.Subjectoid,
		Vid: op.Subjectvid,
	})
	if err != nil {
		return fmt.Errorf("SelectObjective/src: %w", err)
	}

	_, err = q.UpdateActiveVidForObjective(ctx, database.UpdateActiveVidForObjectiveParams{
		Oid: op.Subjectoid,
		Vid: params.To,
	})
	if err != nil {
		return fmt.Errorf("UpdateActiveVidForObjective: %w", err)
	}

	deltas, err := a.calculateDeltasForTwoVersions(ctx, q, srcobj.Bupid, dstobj.Bupid)
	if err != nil {
		return fmt.Errorf("calculateDeltasForTwoVersions: %w", err)
	}

	_, err = a.bubblink(ctx, q, slices.Insert(activepath, 0, models.Ovid{dstobj.Oid, dstobj.Vid}), op, deltas)
	if err != nil {
		return fmt.Errorf("bubblink: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
