-- name: SelectCollaborationOnControlArea :one
SELECT
    *
FROM
    "collaboration"
WHERE
    "aid" = $1
    AND "deleted_at" IS NULL
LIMIT 50;

-- name: SelectCollaboration :one
SELECT
    *
FROM
    "collaboration"
WHERE
    "cid" = $1;

-- name: InsertCollaboration :one
INSERT INTO "collaboration"("cid", "creator", "admin", "leader")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: DeleteCollaboration :exec
UPDATE
    "collaboration"
SET
    "deleted_at" = CURRENT_TIMESTAMP
WHERE
    "cid" = $1;

-- name: InsertCollaborator :one
INSERT INTO "collaborator"("cid", "uid")
    VALUES ($1, $2)
RETURNING
    *;

-- name: SelectCollaborators :many
SELECT
    *
FROM
    "collaborator"
WHERE
    "cid" = $1
    AND "deleted_at" IS NULL
LIMIT 100;

-- name: DeleteCollaborator :exec
UPDATE
    "collaborator"
SET
    "deleted_at" = CURRENT_TIMESTAMP
WHERE
    "cid" = $1
    AND "uid" = $2;

