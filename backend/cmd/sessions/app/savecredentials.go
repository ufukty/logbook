package app

import (
	"context"
	"fmt"
	"logbook/cmd/sessions/database"
	"logbook/models/columns"
	"logbook/models/transports"

	"github.com/alexedwards/argon2id"
)

var argon2idParams = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

type SaveCredentialsRequest struct {
	Uid      columns.UserId      `json:"uid"`
	Email    columns.Email       `json:"email"`
	Password transports.Password `json:"password"`
}

var ErrEmailExists = fmt.Errorf("email in use")

func (a *App) SaveCredentials(ctx context.Context, params SaveCredentialsRequest) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	_, err = q.SelectLatestLoginByEmail(ctx, params.Email)
	if err == nil {
		return ErrEmailExists
	}

	hash, err := argon2id.CreateHash(string(params.Password), argon2idParams)
	if err != nil {
		return fmt.Errorf("generating hash: %w", err)
	}

	_, err = q.InsertLogin(ctx, database.InsertLoginParams{
		Uid:   params.Uid,
		Email: params.Email,
		Hash:  hash,
	})
	if err != nil {
		return fmt.Errorf("inserting login information into database: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
