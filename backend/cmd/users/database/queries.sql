-- name: InsertUser :one
INSERT INTO "user" DEFAULT
    VALUES
    RETURNING
        *;

-- name: SelectUserByUserId :one
SELECT
    *
FROM
    "user"
WHERE
    "uid" = $1
LIMIT 1;

-- name: DeleteUserByUid :exec
UPDATE
    "user"
SET
    "deleted" = TRUE
WHERE
    "uid" = $1;

