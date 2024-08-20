-- name: SelectActive :one
SELECT
    *
FROM
    "active"
WHERE
    "oid" = $1;

-- name: InsertActiveVidForObjective :one
INSERT INTO "active"("oid", "vid")
    VALUES ($1, $2)
RETURNING
    *;

-- name: UpdateActiveVidForObjective :one
UPDATE
    "active"
SET
    "vid" = $2
WHERE
    "oid" = $1
RETURNING
    *;

