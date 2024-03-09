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

;

-- name: InsertLogin :one
INSERT INTO "login"("uid", "email", "hash")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: DeleteLoginByLid :exec
UPDATE
    "login"
SET
    "deleted" = TRUE
WHERE
    "lid" = $1;

-- name: SelectLoginsByUid :many
SELECT
    *
FROM
    "login"
WHERE
    "uid" = $1
    AND ! "deleted";

-- name: SelectLatestLoginByEmail :one
SELECT
    *
FROM
    "login"
WHERE
    "email" = $1
    AND ! "deleted"
ORDER BY
    "created_at"
LIMIT 1;

;

-- name: InsertSession :one
INSERT INTO "access"("uid", "uid")
    VALUES ($1, $2)
RETURNING
    *;

;

-- name: InsertAccess :one
INSERT INTO "access"("uid", "useragent", "ipaddress")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: SelectLatestTwentyAccessesByUid :many
SELECT
    *
FROM
    "access"
WHERE
    "uid" = $1
ORDER BY
    "created_at"
LIMIT 20;

