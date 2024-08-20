-- name: SelectProperties :one
SELECT
    *
FROM
    "props"
WHERE
    "pid" = $1
LIMIT 1;

-- name: InsertProperties :one
INSERT INTO "props"("content", "completed", "creator", "owner")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

