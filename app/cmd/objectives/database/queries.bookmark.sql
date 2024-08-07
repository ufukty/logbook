-- name: InsertBookmark :one
INSERT INTO "bookmark"("uid", "oid", "vid", "title", "is_rock")
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: SelectBookmarks :many
SELECT
    *
FROM
    "bookmark"
WHERE
    "uid" == $1
LIMIT 100;

