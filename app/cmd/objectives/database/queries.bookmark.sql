-- name: InsertBookmark :one
INSERT INTO "bookmark"("uid", "oid", "vid", "display_name", "is_rock")
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

