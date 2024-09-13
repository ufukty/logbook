package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
	"logbook/cmd/objectives/endpoints"
	"logbook/models/columns"

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
	Firstname columns.HumanName
	Lastname  columns.HumanName
	Email     columns.Email
	Password  string
}

var ErrEmailExists = fmt.Errorf("email in use")

func (a *App) CreateUser(ctx context.Context, params RegistrationParameters) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	_, err = q.SelectLatestLoginByEmail(ctx, string(params.Email))
	if err == nil {
		return ErrEmailExists
	}

	user, err := q.InsertUser(ctx)
	if err != nil {
		return fmt.Errorf("inserting record to database: %w", err)
	}

	hash, err := argon2id.CreateHash(params.Password, argon2idParams)
	if err != nil {
		return fmt.Errorf("generating hash: %w", err)
	}

	_, err = q.InsertLogin(ctx, database.InsertLoginParams{
		Uid:   user.Uid,
		Email: string(params.Email),
		Hash:  hash,
	})
	if err != nil {
		return fmt.Errorf("inserting login information into database: %w", err)
	}

	_, err = q.InsertProfileInformation(ctx, database.InsertProfileInformationParams{
		Uid:       user.Uid,
		Firstname: string(params.Firstname),
		Lastname:  string(params.Lastname),
	})
	if err != nil {
		return fmt.Errorf("inserting profile information into database: %w", err)
	}

	r, err := a.objectives.RockCreate(&endpoints.RockCreateRequest{
		UserId: user.Uid,
	})
	if err != nil {
		return fmt.Errorf("calling objectives service to create the rock for user: %w", err)
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("objectives service returned non-200 status code for the request to create rock for user: %s", r.Body)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
