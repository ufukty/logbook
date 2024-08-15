-- name: InsertComputedToTop :one
INSERT INTO "computed_to_top"("oid", "vid", "viewer", "viewer_type", "is_solo", "is_completed", "index", "subtree_size", "completed_subitems")
    VALUES ("$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9")
RETURNING
    *;

-- name: SelectComputedToTop :one
SELECT
    *
FROM
    "computed_to_top"
WHERE
    "oid" = $1
    AND "vid" = $2
    AND "viewer" = $3
LIMIT 1;

