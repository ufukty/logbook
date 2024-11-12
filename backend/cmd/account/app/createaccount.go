package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
	"logbook/cmd/objectives/endpoints"
	"logbook/models/columns"
	"time"

	"github.com/alexedwards/argon2id"
)

var argon2idParams = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

type CreateAccountRequest struct {
	CsrfToken string

	Firstname columns.HumanName
	Lastname  columns.HumanName
	Birthday  time.Time
	// Country   columns.Country

	Email columns.Email
	// EmailGrant columns.EmailGrant

	// Phone      columns.Phone
	// PhoneGrant columns.PhoneGrant

	Password string
}

var ErrEmailExists = fmt.Errorf("email in use")

func (a *App) CreateAccount(ctx context.Context, params CreateAccountRequest) error {
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
		Email: params.Email,
		Hash:  hash,
	})
	if err != nil {
		return fmt.Errorf("inserting login information into database: %w", err)
	}

	_, err = q.InsertProfileInformation(ctx, database.InsertProfileInformationParams{
		Uid:       user.Uid,
		Firstname: params.Firstname,
		Lastname:  params.Lastname,
	})
	if err != nil {
		return fmt.Errorf("inserting profile information into database: %w", err)
	}

	_, err = a.objectives.RockCreate(&endpoints.RockCreateRequest{
		UserId: user.Uid,
	})
	if err != nil {
		return fmt.Errorf("creating rock for user via objectives service: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
