-- DROP DATABASE IF EXISTS logbook;
-- CREATE DATABASE logbook;
-- \c logbook;
-- ;
CREATE DOMAIN "UserId" AS uuid;

CREATE DOMAIN "ObjectiveId" AS uuid;

CREATE DOMAIN "VersionId" AS uuid;

CREATE DOMAIN "CommitId" AS uuid;

CREATE DOMAIN "OperationId" AS uuid;

CREATE DOMAIN "LinkId" AS uuid;

CREATE DOMAIN "BookmarkId" AS uuid;

CREATE TABLE "version"(
    "vid" "VersionId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "based" "VersionId" NOT NULL
);

CREATE TABLE "objective"(
    "oid" "ObjectiveId" NOT NULL DEFAULT gen_random_uuid(), -- objective id
    "vid" "VersionId" NOT NULL,
    "based" "VersionId" NOT NULL, -- previous "vid" OR '00000000-0000-0000-0000-000000000000'
    "content" text NOT NULL,
    "creator" "UserId" NOT NULL, -- user id
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("oid", "vid")
);

CREATE INDEX "index_objective" ON "objective"("created_at");

CREATE TABLE "objective_link"(
    "lid" "LinkId" NOT NULL DEFAULT gen_random_uuid(), -- link id
    "sup_oid" "ObjectiveId" NOT NULL, -- super objective id
    "sup_vid" "VersionId" NOT NULL, -- super version id
    "sub_oid" "ObjectiveId" NOT NULL, -- sub objective id
    "sub_vid" "VersionId" NOT NULL, -- sub version id
    "creator" "UserId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("sup_oid", "sup_vid", "sub_oid")
);

CREATE TABLE "objective_completion"(
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL,
    "actor" "UserId" NOT NULL, -- user id
    "completed" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- MARK: Computed properties and user preferences per item per user
;

CREATE TABLE "objective_view"(
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL,
    "viewer" "UserId" NOT NULL,
    "degree" int NOT NULL,
    "depth" int NOT NULL,
    "ready" boolean NOT NULL,
    "completion_pct" float NOT NULL,
    "fold" boolean NOT NULL
);

CREATE TABLE "computed_to_top"(
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL,
    "dependencies_are_cleared" boolean NOT NULL,
    "all_cleared" boolean NOT NULL,
    "degree" int NOT NULL,
    "completed_subtasks" int NOT NULL,
    PRIMARY KEY ("oid", "vid")
);

CREATE TABLE "computed_to_bottom"(
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL,
    "depth" int NOT NULL,
    PRIMARY KEY ("oid", "vid")
);

-- MARK: operations
;

CREATE TABLE "op_objective_create"(
    "opid" "OperationId" NOT NULL DEFAULT gen_random_uuid(),
    "actor" "UserId" NOT NULL,
    "poid" "ObjectiveId" NOT NULL, -- parent oid
    "pvid" "VersionId" NOT NULL, -- parent vid based on
    "content" text,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_delete"(
    "opid" "OperationId" NOT NULL DEFAULT gen_random_uuid(),
    "actor" "UserId" NOT NULL,
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL, -- based on
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_content_update"(
    "opid" "OperationId" NOT NULL DEFAULT gen_random_uuid(),
    "actor" "UserId" NOT NULL,
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL, -- based on
    "content" text,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_attach_subobjective"(
    "opid" "OperationId" NOT NULL DEFAULT gen_random_uuid(),
    "actor" "UserId" NOT NULL,
    "sup_oid" "ObjectiveId" NOT NULL,
    "sup_vid" "VersionId" NOT NULL, -- based on
    "sub_oid" "ObjectiveId" NOT NULL,
    "sub_vid" "VersionId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_update_completion"(
    "opid" "OperationId" NOT NULL DEFAULT gen_random_uuid(),
    "actor" "UserId" NOT NULL,
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL,
    "completed" boolean NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- TODO: operations: reorder, note (create,update,delete), delegation (assign, unassign), collaboration (init, invite, restrict), versioning (rollback, fastforward)
;

CREATE TABLE "bookmark"(
    "bid" "BookmarkId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "oid" "ObjectiveId" NOT NULL,
    "vid" "VersionId" NOT NULL,
    "display_name" text,
    "is_rock" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

