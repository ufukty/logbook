-- name: InsertComputedToTop :one
INSERT INTO "computed_to_top"("oid", "vid", "is_in_collaboration")
    VALUES ("$1", "$2", "$3")
RETURNING
    *;

-- name: UpdateComputedToTop :one
UPDATE
    "computed_to_top"
SET
    "is_in_collaboration" = $3
WHERE
    "oid" = $1
    AND "vid" = $2
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
LIMIT 1;

-- name: InsertComputedToTopSolo :one
INSERT INTO "computed_to_top_solo"("oid", "vid", "is_completed", "subtree_size", "completed_subitems")
    VALUES ("$1", "$2", "$3", "$4", "$5")
RETURNING
    *;

-- name: SelectComputedToTopSolo :one
SELECT
    *
FROM
    "computed_to_top_solo"
WHERE
    "oid" = $1
    AND "vid" = $2
LIMIT 1;

-- name: InsertComputedToTopCollaborated :one
INSERT INTO "computed_to_top_collaborated"("oid", "vid", "is_completed", "subtree_size", "completed_subitems")
    VALUES ("$1", "$2", "$3", "$4", "$5")
RETURNING
    *;

-- name: SelectComputedToTopCollaborated :one
SELECT
    *
FROM
    "computed_to_top_collaborated"
WHERE
    "oid" = $1
    AND "vid" = $2
LIMIT 1;

-- name: InsertComputedToTopCollaborator :one
INSERT INTO "computed_to_top_collaborator"("oid", "vid", "viewer", "subtree_size", "completed_subitems")
    VALUES ("$1", "$2", "$3", "$4", "$5")
RETURNING
    *;

-- name: SelectComputedToTopCollaborator :one
SELECT
    *
FROM
    "computed_to_top_collaborator"
WHERE
    "oid" = $1
    AND "vid" = $2
    AND "viewer" = $3
LIMIT 1;

