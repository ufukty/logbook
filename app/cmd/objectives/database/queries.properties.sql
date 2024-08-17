-- name: SelectProperties :one
SELECT
    *
FROM
    "props"
WHERE
    "pid" = $1
LIMIT 1;

-- name: InsertProperties :one
INSERT INTO "props"("content", "creator")
    VALUES ($1, $2)
RETURNING
    *;

