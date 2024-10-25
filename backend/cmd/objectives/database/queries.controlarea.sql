-- name: InsertControlArea :one
INSERT INTO "control_area"("root", "catype")
    VALUES ($1, $2)
RETURNING
    *;

-- name: DeleteControlArea :one
UPDATE
    "control_area"
SET
    "deleted_at" = TIMESTAMP
WHERE
    "caid" = $1
RETURNING
    *;

-- name: SelectControlAreaOnObjective :one
SELECT
    *
FROM
    "control_area"
WHERE
    "root" = $1
LIMIT 1;

