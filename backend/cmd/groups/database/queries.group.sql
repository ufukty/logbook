-- name: InsertNewGroup :one
INSERT INTO "group"("name", "creator")
    VALUES ($1, $2)
RETURNING
    *;

-- name: UpdateGroupName :one
UPDATE
    "group"
SET
    "name" = $2
WHERE
    "gid" = $1
RETURNING
    *;

-- name: SelectGroupsCreatedBy :many
SELECT
    *
FROM
    "group"
WHERE
    "creator" = $1
    AND "deleted_at" IS NOT NULL -- FIXME:
LIMIT 20;

-- name: SelectUserTypeGroupMembers :many
SELECT
    *
FROM
    "group_member_user"
WHERE
    "gid" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200;

-- name: SelectGroupTypeGroupMembers :many
SELECT
    *
FROM
    "group_member_group"
WHERE
    "gid" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200;

-- name: SelectGroupsByUserTypeMember :many
SELECT
    *
FROM
    "group_member_user"
WHERE
    "member" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200;

-- name: SelectGroupsByGroupTypeMember :many
SELECT
    *
FROM
    "group_member_group"
WHERE
    "member" = $1
    AND "deleted_at" IS NOT NULL
LIMIT 200;

