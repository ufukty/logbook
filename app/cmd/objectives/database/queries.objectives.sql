-- name: SelectObjective :one
SELECT
    *
FROM
    "objective"
WHERE
    "oid" = $1
    AND "vid" = $2
LIMIT 1;

-- name: InsertNewObjective :one
INSERT INTO "objective"("based", "created_by", "props")
    VALUES ('00000000-0000-0000-0000-000000000000', $1, $2)
RETURNING
    *;

-- name: InsertUpdatedObjective :one
INSERT INTO "objective"("oid", "based", "created_by", "props")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: InsertRock :one
INSERT INTO "objective"("oid", "based", "created_by", "props")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

