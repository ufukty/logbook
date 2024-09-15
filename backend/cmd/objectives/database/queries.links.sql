-- name: SelectSubLinks :many
SELECT
    *
FROM
    "link"
WHERE
    "sup_oid" = $1
    AND "sup_vid" = $2
ORDER BY
    "created_at_original" ASC
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

-- name: SelectUpperLinksToActiveObjectives :many
SELECT
    "sup_oid",
    "sup_vid"
FROM
    "link" AS l
    INNER JOIN "active" AS a ON l."sup_oid" = a."oid"
        AND l."sup_vid" = a."vid"
WHERE
    l."sub_oid" = $1
    AND l."sub_vid" = $2;

-- name: InsertNewLink :one
INSERT INTO "link"("sup_oid", "sup_vid", "sub_oid", "sub_vid")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: InsertUpdatedLink :one
INSERT INTO "link"("sup_oid", "sup_vid", "sub_oid", "sub_vid", "created_at_original")
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

