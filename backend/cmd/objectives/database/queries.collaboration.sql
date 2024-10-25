-- name: SelectCollaborationOnControlArea :one
SELECT
    *
FROM
    "collaboration"
WHERE
    "caid" = $1
    AND "deleted_at" IS NULL
LIMIT 50;

-- name: SelectCollaboration :one
SELECT
    *
FROM
    "collaboration"
WHERE
    "coid" = $1;

-- name: InsertCollaboration :one
INSERT INTO "collaboration"("coid", "creator", "admin", "leader")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: DeleteCollaboration :exec
UPDATE
    "collaboration"
SET
    "deleted_at" = CURRENT_TIMESTAMP
WHERE
    "coid" = $1;

-- name: InsertGroupTypeCollaborator :one
INSERT INTO "collaborator_group"("coid", "gid")
    VALUES ($1, $2)
RETURNING
    *;

-- name: ListGroupTypeCollaborators :many
SELECT
    *
FROM
    "collaborator_group"
WHERE
    "coid" = $1
    AND "deleted_at" IS NULL
LIMIT 100;

-- name: SelectGroupTypeCollaboratorByGroupId :one
SELECT
    *
FROM
    "collaborator_group"
WHERE
    "coid" = $1
    AND "gid" = $2
    AND "deleted_at" IS NULL
LIMIT 1;

-- name: DeleteGroupTypeCollaboratorByCollaboratorId :exec
UPDATE
    "collaborator_group"
SET
    "deleted_at" = CURRENT_TIMESTAMP
WHERE
    "coid" = $1
    AND "crid" = $2;

-- name: InsertUserTypeCollaborator :one
INSERT INTO "collaborator_user"("coid", "uid")
    VALUES ($1, $2)
RETURNING
    *;

-- name: ListUserTypeCollaborators :many
SELECT
    *
FROM
    "collaborator_user"
WHERE
    "coid" = $1
    AND "deleted_at" IS NULL
LIMIT 100;

-- name: SelectUserTypeCollaboratorByUserId :one
SELECT
    *
FROM
    "collaborator_user"
WHERE
    "coid" = $1
    AND "uid" = $2
    AND "deleted_at" IS NULL
LIMIT 1;

-- name: DeleteUserTypeCollaboratorByCollaboratorId :exec
UPDATE
    "collaborator_user"
SET
    "deleted_at" = CURRENT_TIMESTAMP
WHERE
    "coid" = $1
    AND "crid" = $2;

