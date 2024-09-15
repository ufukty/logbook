-- name: SelectTagsByObjectiveId :many
SELECT
  *
FROM
  "tagging"
WHERE
  "oid" = $1
  AND "vid" = $2
  AND NOT "deleted";

