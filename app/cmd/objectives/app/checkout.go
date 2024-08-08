package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5"
)

type CheckoutParams struct {
	User    columns.UserId
	Subject models.Ovid
	To      columns.VersionId
}

var ErrVersionDoesNotExist = fmt.Errorf("given version of the objective doesn't exist")

// TODO: bubblink (update and checkout ascendants) (implicit/explicit checkouts?)
func (a *App) Checkout(ctx context.Context, params CheckoutParams) error {
	activepath, err := a.listActivePathToRock(ctx, params.Subject)
	if err == ErrLeftBehind {
		return ErrLeftBehind
	} else if err != nil {
		return fmt.Errorf("checking if the objective %s is in active path: %w", params.Subject, err)
	}

	op, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.Subject.Oid,
		Subjectvid: params.Subject.Vid,
		Actor:      params.User,
		OpType:     database.OpTypeCheckout,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("insert checkout operation: %w", err)
	}

	_, err = a.queries.InsertOpCheckout(ctx, database.InsertOpCheckoutParams{
		Opid: op.Opid,
		To:   params.To,
	})
	if err != nil {
		return fmt.Errorf("insert checkout operation details: %w", err)
	}

	obj, err := a.queries.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: op.Subjectoid,
		Vid: params.To,
	})
	if err == pgx.ErrNoRows {
		return ErrVersionDoesNotExist
	} else if err != nil {
		return fmt.Errorf("checking if the version of objective requested exist: %w", err)
	}

	_, err = a.queries.UpdateActiveVidForObjective(ctx, database.UpdateActiveVidForObjectiveParams{
		Oid: op.Subjectoid,
		Vid: params.To,
	})
	if err != nil {
		return fmt.Errorf("updating the active version of objective: %w", err)
	}

	_, err = a.bubblink(ctx, append(activepath, models.Ovid{obj.Oid, obj.Vid}), op)
	if err != nil {
		return fmt.Errorf("promoting the version change to ascendants: %w", err)
	}

	return nil
}
