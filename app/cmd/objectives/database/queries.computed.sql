CREATE TABLE "computed_to_top"(
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL,
    "viewer" "UserId" NOT NULL,
    "is_solo" boolean NOT NULL,
    "is_completed" boolean NOT NULL,
    "index" int NOT NULL,
    "subtree_size" int NOT NULL,
    "completed_subitems" int NOT NULL,
    PRIMARY KEY ("oid", "vid", "viewer")
);

-- name: InsertComputedToTop :one
INSERT INTO "computed_to_top"("oid", "vid", "viewer", "is_solo", "is_completed", "index", "subtree_size", "completed_subitems")
    VALUES ("$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8")
RETURNING
    *;

-- name: SelectComputedToTop :one
SELECT
    *
FROM
    "computed_to_top"
WHERE
    "oid" = $1
    AND "vid" = $2
    AND "viewer" = $3
LIMIT 1;

