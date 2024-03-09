package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"

	"github.com/alexedwards/argon2id"
)

var argon2idParams = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

type RegistrationParameters struct {
	Username string
	Email    string
	Password string
}

func (a *App) Register(ctx context.Context, params RegistrationParameters) error {
	user, err := a.queries.InsertUser(ctx)
	if err != nil {
		return fmt.Errorf("inserting record to database: %w", err)
	}

	hash, err := argon2id.CreateHash(params.Password, argon2idParams)
	if err != nil {
		return fmt.Errorf("generating hash: %w", err)
	}

	_, err = a.queries.InsertLogin(ctx, database.InsertLoginParams{
		Uid:   user.Uid,
		Email: params.Email,
		Hash:  hash,
	})
	if err != nil {
		return fmt.Errorf("inserting login information into database: %w", err)
	}

	return nil
}
