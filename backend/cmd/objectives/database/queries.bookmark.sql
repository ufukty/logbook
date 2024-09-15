-- name: InsertBookmark :one
INSERT INTO "bookmark"("uid", "oid", "title", "is_rock")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: SelectBookmarks :many
SELECT
    *
FROM
    "bookmark"
WHERE
    "uid" = $1
LIMIT 100;

-- name: SelectTheRockForUser :one
SELECT
    *
FROM
    "bookmark"
WHERE
    "uid" = $1
    AND "is_rock" = TRUE
LIMIT 1;

