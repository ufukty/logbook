-- name: InsertProfileInformation :one
INSERT INTO "profile"("uid", "firstname", "lastname")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: SelectLatestProfileByUid :one
SELECT
    *
FROM
    "profile"
WHERE
    "uid" = $1
ORDER BY
    "created_at" DESC
LIMIT 1;

