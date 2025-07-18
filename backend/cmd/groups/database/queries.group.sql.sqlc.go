// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries.group.sql

package database

import (
	"context"

	"logbook/models/columns"
)

const insertNewGroup = `-- name: InsertNewGroup :one
INSERT INTO "group"("name", "creator")
    VALUES ($1, $2)
RETURNING
    gid, name, creator, created_at, deleted_at
`

type InsertNewGroupParams struct {
	Name    string
	Creator columns.UserId
}

func (q *Queries) InsertNewGroup(ctx context.Context, arg InsertNewGroupParams) (Group, error) {
	row := q.db.QueryRow(ctx, insertNewGroup, arg.Name, arg.Creator)
	var i Group
	err := row.Scan(
		&i.Gid,
		&i.Name,
		&i.Creator,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const selectGroupTypeGroupMembers = `-- name: SelectGroupTypeGroupMembers :many
SELECT
    gmid, gid, member, ginvid, created_at, deleted_at
FROM
    "group_member_group"
WHERE
    "gid" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200
`

func (q *Queries) SelectGroupTypeGroupMembers(ctx context.Context, gid columns.GroupId) ([]GroupMemberGroup, error) {
	rows, err := q.db.Query(ctx, selectGroupTypeGroupMembers, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupMemberGroup
	for rows.Next() {
		var i GroupMemberGroup
		if err := rows.Scan(
			&i.Gmid,
			&i.Gid,
			&i.Member,
			&i.Ginvid,
			&i.CreatedAt,
			&i.DeletedAt,
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

const selectGroupsByGroupTypeMember = `-- name: SelectGroupsByGroupTypeMember :many
SELECT
    gmid, gid, member, ginvid, created_at, deleted_at
FROM
    "group_member_group"
WHERE
    "member" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200
`

func (q *Queries) SelectGroupsByGroupTypeMember(ctx context.Context, member columns.GroupId) ([]GroupMemberGroup, error) {
	rows, err := q.db.Query(ctx, selectGroupsByGroupTypeMember, member)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupMemberGroup
	for rows.Next() {
		var i GroupMemberGroup
		if err := rows.Scan(
			&i.Gmid,
			&i.Gid,
			&i.Member,
			&i.Ginvid,
			&i.CreatedAt,
			&i.DeletedAt,
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

const selectGroupsByUserTypeMember = `-- name: SelectGroupsByUserTypeMember :many
SELECT
    gmid, gid, member, ginvid, created_at, deleted_at
FROM
    "group_member_user"
WHERE
    "member" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200
`

func (q *Queries) SelectGroupsByUserTypeMember(ctx context.Context, member columns.UserId) ([]GroupMemberUser, error) {
	rows, err := q.db.Query(ctx, selectGroupsByUserTypeMember, member)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupMemberUser
	for rows.Next() {
		var i GroupMemberUser
		if err := rows.Scan(
			&i.Gmid,
			&i.Gid,
			&i.Member,
			&i.Ginvid,
			&i.CreatedAt,
			&i.DeletedAt,
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

const selectGroupsCreatedBy = `-- name: SelectGroupsCreatedBy :many
SELECT
    gid, name, creator, created_at, deleted_at
FROM
    "group"
WHERE
    "creator" = $1
    AND "deleted_at" IS NOT NULL -- FIXME:
LIMIT 20
`

func (q *Queries) SelectGroupsCreatedBy(ctx context.Context, creator columns.UserId) ([]Group, error) {
	rows, err := q.db.Query(ctx, selectGroupsCreatedBy, creator)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Group
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.Gid,
			&i.Name,
			&i.Creator,
			&i.CreatedAt,
			&i.DeletedAt,
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

const selectUserTypeGroupMembers = `-- name: SelectUserTypeGroupMembers :many
SELECT
    gmid, gid, member, ginvid, created_at, deleted_at
FROM
    "group_member_user"
WHERE
    "gid" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200
`

func (q *Queries) SelectUserTypeGroupMembers(ctx context.Context, gid columns.GroupId) ([]GroupMemberUser, error) {
	rows, err := q.db.Query(ctx, selectUserTypeGroupMembers, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupMemberUser
	for rows.Next() {
		var i GroupMemberUser
		if err := rows.Scan(
			&i.Gmid,
			&i.Gid,
			&i.Member,
			&i.Ginvid,
			&i.CreatedAt,
			&i.DeletedAt,
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

const updateGroupName = `-- name: UpdateGroupName :one
UPDATE
    "group"
SET
    "name" = $2
WHERE
    "gid" = $1
RETURNING
    gid, name, creator, created_at, deleted_at
`

type UpdateGroupNameParams struct {
	Gid  columns.GroupId
	Name string
}

func (q *Queries) UpdateGroupName(ctx context.Context, arg UpdateGroupNameParams) (Group, error) {
	row := q.db.QueryRow(ctx, updateGroupName, arg.Gid, arg.Name)
	var i Group
	err := row.Scan(
		&i.Gid,
		&i.Name,
		&i.Creator,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}
