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
    AND NOT "deleted";

-- name: SelectLatestLoginByEmail :one
SELECT
    *
FROM
    "login"
WHERE
    "email" = $1
    AND NOT "deleted"
ORDER BY
    "created_at"
LIMIT 1;

;

-- name: InsertProfileInformation :one
INSERT INTO "profile"("uid", "firstname", "lastname")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: SelectProfileByUid :one
SELECT
    *
FROM
    "profile"
WHERE
    "uid" = $1
ORDER BY
    "created_at" DESC
LIMIT 1;

;

-- name: InsertSession :one
INSERT INTO "session_standard"("uid", "token")
    VALUES ($1, $2)
RETURNING
    *;

-- name: SelectSessionByToken :one
SELECT
    *
FROM
    "session_standard"
WHERE
    "token" = $1;

-- name: SelectActiveSessionsByUid :many
SELECT
    *
FROM
    "session_standard"
WHERE
    "uid" = $1
    AND NOT "deleted";

-- name: DeleteSessionBySid :exec
UPDATE
    "session_standard"
SET
    "deleted" = TRUE
WHERE
    "sid" = $1;

;

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

;

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

