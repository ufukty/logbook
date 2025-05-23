-- name: SelectOperation :one
SELECT
    *
FROM
    "operation"
WHERE
    "opid" = $1
LIMIT 1;

-- name: InsertOperation :one
INSERT INTO "operation"("subjectoid", "subjectvid", "actor", "op_type", "op_status")
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: InsertOpCheckout :one
INSERT INTO "op_checkout"("opid", "to")
    VALUES ($1, $2)
RETURNING
    *;

-- name: InsertOpObjCompletion :one
INSERT INTO "op_obj_completion"("opid", "completed")
    VALUES ($1, $2)
RETURNING
    *;

-- name: InsertOpObjContent :one
INSERT INTO "op_obj_content"("opid", "content")
    VALUES ($1, $2)
RETURNING
    *;

-- name: InsertOpObjInit :one
INSERT INTO "op_obj_init"("opid", "content")
    VALUES ($1, $2)
RETURNING
    *;

-- name: InsertOpObjCreateSubtask :one
INSERT INTO "op_obj_create_subtask"("opid", "soid", "svid")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: InsertOpObjDeleteSubtask :one
INSERT INTO "op_obj_delete_subtask"("opid", "doid", "dvid")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: InsertOpObjAttach :one
INSERT INTO "op_obj_attach"("opid", "child")
    VALUES ($1, $2)
RETURNING
    *;

-- name: InsertOpObjDetach :one
INSERT INTO "op_obj_detach"("opid", "child")
    VALUES ($1, $2)
RETURNING
    *;

-- name: InsertOpObjReorder :one
INSERT INTO "op_obj_reorder"("opid", "child", "moveafter")
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: InsertOpTransitive :one
INSERT INTO "op_transitive"("opid", "cause")
    VALUES ($1, $2)
RETURNING
    *;

-- name: InsertOpDoubleTransitiveMerger :one
INSERT INTO "op_double_transitive_merger"("opid", "first", "second")
    VALUES ($1, $2, $3)
RETURNING
    *;

