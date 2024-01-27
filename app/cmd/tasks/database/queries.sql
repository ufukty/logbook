-- MARK: objectives
;

-- name: SelectObjective :one
SELECT
    "oid",
    "vid",
    "based",
    "content",
    "creator",
    "created_at"
FROM
    "objective"
WHERE
    "oid" = $1
    AND "vid" = $2
LIMIT 1;

-- name: InsertObjective :one
INSERT INTO "objective"("vid", "based", "content", "creator")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: SelectEffectiveVersionOfObjective :one
SELECT
    "vid"
FROM
    "objective_effective_version"
WHERE
    "oid" = $1
LIMIT 1;

-- name: CreateTask :one
INSERT INTO "objective"("based", "content", "creator")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- MARK: links
;

-- name: SelectSubLinks :many
SELECT
    "lid",
    "sup_oid",
    "sup_vid",
    "sub_oid",
    "sub_vid",
    "created_at"
FROM
    "objective_link"
WHERE
    "sup_oid" = $1
    AND "sup_vid" = $2
LIMIT 50;

-- name: SelectTheUpperLink :one
SELECT
    "lid",
    "sup_oid",
    "sup_vid",
    "sub_oid",
    "sub_vid",
    "created_at"
FROM
    "objective_link"
WHERE
    "sub_oid" = $1
    AND "sub_vid" = $2
LIMIT 1;

-- name: InsertLink :one
INSERT INTO "objective_link"("sup_oid", "sup_vid", "sub_oid", "sub_vid", "creator")
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- MARK: versioning
;

-- name: SelectVersion :one
SELECT
    "vid",
    "based"
FROM
    "version"
WHERE
    "vid" = $1
LIMIT 1;

-- name: InsertVersion :one
INSERT INTO "version"("based")
    VALUES ($1)
RETURNING
    *;

-- MARK: versioning config
;

-- name: SelectVersioningConfig :one
SELECT
    "oid",
    "first",
    "effective"
FROM
    "versioning_config"
WHERE
    "oid" = $1
LIMIT 1;

;

-- MARK: operations
;

-- name: InsertOpObjectiveCreate :one
INSERT INTO op_objective_create("poid", "pvid", "actor", "content")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: InsertOpObjectiveDelete :one
INSERT INTO op_objective_delete("oid", "vid", "actor")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: InsertOpObjectiveContentUpdate :one
INSERT INTO op_objective_content_update("oid", "vid", "actor", "content")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: InsertOpObjectiveAttachSubobjective :one
INSERT INTO op_objective_attach_subobjective("actor", "sup_oid", "sup_vid", "sub_oid", "sub_vid")
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: InsertOpObjectiveUpdateCompletion :one
INSERT INTO op_objective_update_completion("oid", "vid", "actor", "completed")
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

