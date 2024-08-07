-- name: InsertObjectiveViewPrefs :one
INSERT INTO "objective_view_prefs"("uid", "oid", "fold")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: SelectObjectiveViewPrefs :one
SELECT
    *
FROM
    "objective_view_prefs"
WHERE
    "uid" = $1
    AND "oid" = $2
LIMIT 1;

