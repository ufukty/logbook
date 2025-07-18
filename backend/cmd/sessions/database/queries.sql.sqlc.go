// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries.sql

package database

import (
	"context"
	"net/netip"

	"logbook/models/columns"
)

const deleteLoginByLid = `-- name: DeleteLoginByLid :exec
UPDATE
    "login"
SET
    "deleted" = TRUE
WHERE
    "lid" = $1
`

func (q *Queries) DeleteLoginByLid(ctx context.Context, lid columns.LoginId) error {
	_, err := q.db.Exec(ctx, deleteLoginByLid, lid)
	return err
}

const deleteSessionBySid = `-- name: DeleteSessionBySid :exec
UPDATE
    "session_standard"
SET
    "deleted" = TRUE
WHERE
    "sid" = $1
`

func (q *Queries) DeleteSessionBySid(ctx context.Context, sid columns.SessionId) error {
	_, err := q.db.Exec(ctx, deleteSessionBySid, sid)
	return err
}

const deleteSessionByToken = `-- name: DeleteSessionByToken :exec
UPDATE
    "session_standard"
SET
    "deleted" = TRUE
WHERE
    "token" = $1
`

func (q *Queries) DeleteSessionByToken(ctx context.Context, token columns.SessionToken) error {
	_, err := q.db.Exec(ctx, deleteSessionByToken, token)
	return err
}

const insertAccess = `-- name: InsertAccess :one
INSERT INTO "access"("uid", "useragent", "ipaddress")
    VALUES ($1, $2, $3)
RETURNING
    aid, uid, useragent, ipaddress, created_at
`

type InsertAccessParams struct {
	Uid       columns.UserId
	Useragent columns.UserAgent
	Ipaddress netip.Addr
}

func (q *Queries) InsertAccess(ctx context.Context, arg InsertAccessParams) (Access, error) {
	row := q.db.QueryRow(ctx, insertAccess, arg.Uid, arg.Useragent, arg.Ipaddress)
	var i Access
	err := row.Scan(
		&i.Aid,
		&i.Uid,
		&i.Useragent,
		&i.Ipaddress,
		&i.CreatedAt,
	)
	return i, err
}

const insertLogin = `-- name: InsertLogin :one
INSERT INTO "login"("uid", "email", "hash")
    VALUES ($1, $2, $3)
RETURNING
    lid, uid, email, hash, deleted, created_at
`

type InsertLoginParams struct {
	Uid   columns.UserId
	Email columns.Email
	Hash  string
}

func (q *Queries) InsertLogin(ctx context.Context, arg InsertLoginParams) (Login, error) {
	row := q.db.QueryRow(ctx, insertLogin, arg.Uid, arg.Email, arg.Hash)
	var i Login
	err := row.Scan(
		&i.Lid,
		&i.Uid,
		&i.Email,
		&i.Hash,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}

const insertSession = `-- name: InsertSession :one
INSERT INTO "session_standard"("uid", "token")
    VALUES ($1, $2)
RETURNING
    sid, uid, token, deleted, created_at
`

type InsertSessionParams struct {
	Uid   columns.UserId
	Token columns.SessionToken
}

func (q *Queries) InsertSession(ctx context.Context, arg InsertSessionParams) (SessionStandard, error) {
	row := q.db.QueryRow(ctx, insertSession, arg.Uid, arg.Token)
	var i SessionStandard
	err := row.Scan(
		&i.Sid,
		&i.Uid,
		&i.Token,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}

const insertSessionAccountRead = `-- name: InsertSessionAccountRead :one
INSERT INTO "session_account_read"("uid", "token")
    VALUES ($1, $2)
RETURNING
    sid, uid, token, deleted, created_at
`

type InsertSessionAccountReadParams struct {
	Uid   columns.UserId
	Token columns.SessionToken
}

func (q *Queries) InsertSessionAccountRead(ctx context.Context, arg InsertSessionAccountReadParams) (SessionAccountRead, error) {
	row := q.db.QueryRow(ctx, insertSessionAccountRead, arg.Uid, arg.Token)
	var i SessionAccountRead
	err := row.Scan(
		&i.Sid,
		&i.Uid,
		&i.Token,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}

const insertSessionAccountWrite = `-- name: InsertSessionAccountWrite :one
INSERT INTO "session_account_write"("uid", "token")
    VALUES ($1, $2)
RETURNING
    sid, uid, token, deleted, created_at
`

type InsertSessionAccountWriteParams struct {
	Uid   columns.UserId
	Token columns.SessionToken
}

func (q *Queries) InsertSessionAccountWrite(ctx context.Context, arg InsertSessionAccountWriteParams) (SessionAccountWrite, error) {
	row := q.db.QueryRow(ctx, insertSessionAccountWrite, arg.Uid, arg.Token)
	var i SessionAccountWrite
	err := row.Scan(
		&i.Sid,
		&i.Uid,
		&i.Token,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}

const selectActiveSessionsByUid = `-- name: SelectActiveSessionsByUid :many
SELECT
    sid, uid, token, deleted, created_at
FROM
    "session_standard"
WHERE
    "uid" = $1
    AND NOT "deleted"
`

func (q *Queries) SelectActiveSessionsByUid(ctx context.Context, uid columns.UserId) ([]SessionStandard, error) {
	rows, err := q.db.Query(ctx, selectActiveSessionsByUid, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SessionStandard
	for rows.Next() {
		var i SessionStandard
		if err := rows.Scan(
			&i.Sid,
			&i.Uid,
			&i.Token,
			&i.Deleted,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectLatestLoginByEmail = `-- name: SelectLatestLoginByEmail :one
SELECT
    lid, uid, email, hash, deleted, created_at
FROM
    "login"
WHERE
    "email" = $1
    AND NOT "deleted"
ORDER BY
    "created_at"
LIMIT 1
`

func (q *Queries) SelectLatestLoginByEmail(ctx context.Context, email columns.Email) (Login, error) {
	row := q.db.QueryRow(ctx, selectLatestLoginByEmail, email)
	var i Login
	err := row.Scan(
		&i.Lid,
		&i.Uid,
		&i.Email,
		&i.Hash,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}

const selectLatestTwentyAccessesByUid = `-- name: SelectLatestTwentyAccessesByUid :many
SELECT
    aid, uid, useragent, ipaddress, created_at
FROM
    "access"
WHERE
    "uid" = $1
ORDER BY
    "created_at"
LIMIT 20
`

func (q *Queries) SelectLatestTwentyAccessesByUid(ctx context.Context, uid columns.UserId) ([]Access, error) {
	rows, err := q.db.Query(ctx, selectLatestTwentyAccessesByUid, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Access
	for rows.Next() {
		var i Access
		if err := rows.Scan(
			&i.Aid,
			&i.Uid,
			&i.Useragent,
			&i.Ipaddress,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectLoginsByUid = `-- name: SelectLoginsByUid :many
SELECT
    lid, uid, email, hash, deleted, created_at
FROM
    "login"
WHERE
    "uid" = $1
    AND NOT "deleted"
`

func (q *Queries) SelectLoginsByUid(ctx context.Context, uid columns.UserId) ([]Login, error) {
	rows, err := q.db.Query(ctx, selectLoginsByUid, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Login
	for rows.Next() {
		var i Login
		if err := rows.Scan(
			&i.Lid,
			&i.Uid,
			&i.Email,
			&i.Hash,
			&i.Deleted,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectSessionAccountRead = `-- name: SelectSessionAccountRead :one
SELECT
    sid, uid, token, deleted, created_at
FROM
    "session_account_read"
WHERE
    "sid" = $1
`

func (q *Queries) SelectSessionAccountRead(ctx context.Context, sid columns.SessionId) (SessionAccountRead, error) {
	row := q.db.QueryRow(ctx, selectSessionAccountRead, sid)
	var i SessionAccountRead
	err := row.Scan(
		&i.Sid,
		&i.Uid,
		&i.Token,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}

const selectSessionAccountWrite = `-- name: SelectSessionAccountWrite :one
SELECT
    sid, uid, token, deleted, created_at
FROM
    "session_account_write"
WHERE
    "sid" = $1
`

func (q *Queries) SelectSessionAccountWrite(ctx context.Context, sid columns.SessionId) (SessionAccountWrite, error) {
	row := q.db.QueryRow(ctx, selectSessionAccountWrite, sid)
	var i SessionAccountWrite
	err := row.Scan(
		&i.Sid,
		&i.Uid,
		&i.Token,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}

const selectSessionByToken = `-- name: SelectSessionByToken :one
SELECT
    sid, uid, token, deleted, created_at
FROM
    "session_standard"
WHERE
    "token" = $1
`

func (q *Queries) SelectSessionByToken(ctx context.Context, token columns.SessionToken) (SessionStandard, error) {
	row := q.db.QueryRow(ctx, selectSessionByToken, token)
	var i SessionStandard
	err := row.Scan(
		&i.Sid,
		&i.Uid,
		&i.Token,
		&i.Deleted,
		&i.CreatedAt,
	)
	return i, err
}
