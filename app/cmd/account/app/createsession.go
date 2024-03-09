package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"logbook/cmd/account/database"
	"time"

	"github.com/alexedwards/argon2id"
)

const month = time.Hour * 24 * 30

type CreateSessionParameters struct {
	Email    string
	Password string
}

var ErrHashMismatch = fmt.Errorf("given password's hash doesn't match with stored hash")

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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

func (a *App) CreateSession(ctx context.Context, params CreateSessionParameters) (database.Session, error) {
	login, err := a.queries.SelectLatestLoginByEmail(ctx, params.Email)
	if err != nil {
		return database.Session{}, fmt.Errorf("selecting latest hash from database: %w", err)
	}

	match, _, err := argon2id.CheckHash(params.Password, login.Hash)
	if err != nil {
		return database.Session{}, fmt.Errorf("comparing hashes: %w", err)
	}

	if !match {
		return database.Session{}, ErrHashMismatch
	}

	if time.Now().Sub(login.CreatedAt.Time) > month {
		if err := renewHash(a.queries, ctx, login, params); err != nil {
			log.Println(fmt.Errorf("could not renew hash for user (uid: %q, old hash date %s): %w\n", login.Uid, login.CreatedAt.Time, err))
		}
	}

	tok, err := generateToken(32)
	if err != nil {
		return database.Session{}, fmt.Errorf("generating session token: %w", err)
	}

	session, err := a.queries.InsertSession(ctx, database.InsertSessionParams{
		Uid:   login.Uid,
		Token: tok,
	})
	if err != nil {
		return database.Session{}, fmt.Errorf("inserting session information into database: %w", err)
	}

	return session, nil
}
