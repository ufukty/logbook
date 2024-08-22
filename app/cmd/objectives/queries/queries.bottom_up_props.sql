-- name: InsertBottomUpProps :one
INSERT INTO "bottom_up_props"("children", "subtree_size", "subtree_completed")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: SelectBottomUpProps :one
SELECT
    *
FROM
    "bottom_up_props"
WHERE
    "bupid" = $1
LIMIT 1;

-- name: InsertBottomUpPropsThirdPerson :one
INSERT INTO "bottom_up_props_third_person"("bupid", "viewer", "children", "subtree_size", "subtree_completed")
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: SelectBottomUpPropsThirdPerson :one
SELECT
    *
FROM
    "bottom_up_props_third_person"
WHERE
    "bupid" = $1
    AND "viewer" = $2
LIMIT 1;

