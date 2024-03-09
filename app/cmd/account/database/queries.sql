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
INSERT INTO "session_standard"("uid", "token")
    VALUES ($1, $2)
RETURNING
    *;

-- name: SelectSession :one
SELECT
    *
FROM
    "session_standard"
WHERE
    "sid" = $1;

-- name: InsertSessionAccountRead :one
INSERT INTO "session_account_read"("uid", "token")
    VALUES ($1, $2)
RETURNING
    *;

-- name: SelectSessionAccountRead :one
SELECT
    *
FROM
    "session_account_read"
WHERE
    "sid" = $1;

-- name: InsertSessionAccountWrite :one
INSERT INTO "session_account_write"("uid", "token")
    VALUES ($1, $2)
RETURNING
    *;

-- name: SelectSessionAccountWrite :one
SELECT
    *
FROM
    "session_account_write"
WHERE
    "sid" = $1;

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

