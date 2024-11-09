package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"logbook/cmd/account/database"
	"logbook/internal/average"
	"logbook/models/columns"
	"time"

	"github.com/alexedwards/argon2id"
)

type CreateSessionParameters struct {
	Email    string
	Password string
}

var ErrHashMismatch = fmt.Errorf("given password's hash doesn't match with stored hash")

func generateToken(length int) (columns.SessionToken, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	b64 := base64.URLEncoding.EncodeToString(bytes)
	cut := string([]byte(b64))[:256]
	return columns.SessionToken(cut), nil
}

func renewHash(q *database.Queries, ctx context.Context, login database.Login, params CreateSessionParameters) error {
	hash, err := argon2id.CreateHash(params.Password, argon2idParams)
	if err != nil {
		return fmt.Errorf("generating hash for password: %w", err)
	}

	err = q.DeleteLoginByLid(ctx, login.Lid)
	if err != nil {
		return fmt.Errorf("deleting old hash from database: %w", err)
	}

	_, err = q.InsertLogin(ctx, database.InsertLoginParams{
		Uid:   login.Uid,
		Email: login.Email,
		Hash:  hash,
	})
	if err != nil {
		return fmt.Errorf("inserting new login information into database: %w", err)
	}

	return nil
}

// TODO: send email
func (a *App) Login(ctx context.Context, params CreateSessionParameters) (database.SessionStandard, error) {
	login, err := a.oneshot.SelectLatestLoginByEmail(ctx, params.Email)
	if err != nil {
		return database.SessionStandard{}, fmt.Errorf("selecting latest hash from database: %w", err)
	}

	match, _, err := argon2id.CheckHash(params.Password, login.Hash)
	if err != nil {
		return database.SessionStandard{}, fmt.Errorf("comparing hashes: %w", err)
	}

	if !match {
		return database.SessionStandard{}, ErrHashMismatch
	}

	if time.Now().Sub(login.CreatedAt.Time) > average.Month {
		if err := renewHash(a.oneshot, ctx, login, params); err != nil {
			return database.SessionStandard{}, fmt.Errorf("could not renew hash for user (uid: %q, old hash date %s): %w", login.Uid, login.CreatedAt.Time, err)
		}
	}

	tok, err := generateToken(256)
	if err != nil {
		return database.SessionStandard{}, fmt.Errorf("generating session token: %w", err)
	}

	session, err := a.oneshot.InsertSession(ctx, database.InsertSessionParams{
		Uid:   login.Uid,
		Token: tok,
	})
	if err != nil {
		return database.SessionStandard{}, fmt.Errorf("inserting session information into database: %w", err)
	}

	return session, nil
}
