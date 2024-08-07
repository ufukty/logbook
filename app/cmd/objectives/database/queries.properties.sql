-- name: SelectProperties :one
SELECT
    *
FROM
    "computed_props"
WHERE
    "propid" = $1
LIMIT 1;

-- name: InsertProperties :one
INSERT INTO "computed_props"("content", "creator")
    VALUES ($1, $2)
RETURNING
    *;

