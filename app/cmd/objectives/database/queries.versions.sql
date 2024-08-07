-- name: SelectActive :one
SELECT
    *
FROM
    "active"
WHERE
    "oid" = $1;

-- name: UpdateActiveVidForObjective :one
UPDATE
    "active"
SET
    "vid" = $2
WHERE
    "oid" = $1
RETURNING
    *;

