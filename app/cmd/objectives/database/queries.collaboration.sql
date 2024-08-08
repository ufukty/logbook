-- name: ListCollaborationsOnObjective :many
SELECT
    *
FROM
    "collaboration"
WHERE
    "oid" = $1
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
INSERT INTO "collaboration"("oid", "creator", "admin", "leader")
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

-- name: ListCollaboratorsForCollaboration :many
SELECT
    *
FROM
    "collaborator"
WHERE
    "cid" = $1
LIMIT 50;

-- name: DeleteCollaboratorFromCollaboration :exec
UPDATE
    "collaborator"
SET
    "deleted_at" = CURRENT_TIMESTAMP
WHERE
    "cid" = $1
    AND "uid" = $2;

