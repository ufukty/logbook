-- name: InsertUser :one
INSERT INTO "user"("uid", "email", "hash")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: InsertAccess :one
INSERT INTO "access"("uid", "useragent", "ipaddress")
    VALUES ($1, $2, $3)
RETURNING
    *;

