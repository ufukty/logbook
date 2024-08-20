-- name: SelectSubLinks :many
SELECT
    *
FROM
    "link"
WHERE
    "sup_oid" = $1
    AND "sup_vid" = $2
LIMIT 50;

-- name: SelectUpperLinks :many
SELECT
    *
FROM
    "link"
WHERE
    "sub_oid" = $1
    AND "sub_vid" = $2
LIMIT 50;

-- name: InsertLink :one
INSERT INTO "link"("sup_oid", "sup_vid", "sub_oid", "sub_vid")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

