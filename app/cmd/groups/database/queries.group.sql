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
    AND "deleted_at" IS NOT NULL
LIMIT 20;

