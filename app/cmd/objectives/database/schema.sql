-- DROP DATABASE IF EXISTS logbook;
-- CREATE DATABASE logbook;
-- \c logbook;
-- ;
CREATE DOMAIN "UserId" AS uuid;

CREATE DOMAIN "ObjectiveId" AS uuid;

CREATE DOMAIN "VersionId" AS uuid;

CREATE DOMAIN "OperationId" AS uuid;

CREATE DOMAIN "PropertyId" AS uuid;

CREATE DOMAIN "LinkId" AS uuid;

CREATE DOMAIN "BookmarkId" AS uuid;

CREATE TABLE "active"(
    "oid" "ObjectiveId" NOT NULL UNIQUE,
    "vid" "VersionId" NOT NULL
);

CREATE TABLE "objective"(
    "oid" "ObjectiveId" NOT NULL DEFAULT gen_random_uuid(), -- objective id
    "vid" "VersionId" NOT NULL,
    "based" "VersionId" NOT NULL, -- previous "vid" OR '00000000-0000-0000-0000-000000000000'
    "created_by" "OperationId" NOT NULL,
    "props" "ProperyId",
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("oid", "vid")
);

CREATE INDEX "index_objective" ON "objective"("created_at");

CREATE TABLE "computed_props"(
    "propid" "PropertyId" NOT NULL DEFAULT gen_random_uuid(),
    "content" text NOT NULL,
    "creator" "UserId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("propid")
);

CREATE TABLE "link"(
    "sup_oid" "ObjectiveId" NOT NULL,
    "sup_vid" "VersionId" NOT NULL,
    "sub_oid" "ObjectiveId" NOT NULL,
    "sub_vid" "VersionId" NOT NULL
);

CREATE TABLE "objective_view_prefs"(
    "uid" "UserId" NOT NULL,
    "oid" "ObjectiveId" NOT NULL,
    "fold" boolean NOT NULL
);

-- is_ready = completed_items == completed_items
-- is_leaf = subtree_size == 0
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

CREATE TYPE "OpType" AS ENUM(
    'checkout',
    'obj_completion',
    'obj_content',
    'obj_create_subtask', --
    'obj_delete', --
    'obj_reattach',
    'obj_reorder',
    'transitive'
);

CREATE TYPE "OpStatus" AS ENUM(
    'received',
    'accepted',
    'rejected'
);

CREATE TABLE "operation"(
    "opid" "OperationId" NOT NULL,
    "subjectoid" "ObjectiveId" NOT NULL,
    "subjectvid" "VersionId" NOT NULL,
    "actor" "UserId" NOT NULL,
    "op_type" "OpType" NOT NULL, -- Transitive, Checkout, Completion, Reattach, Reorder, Content
    "status" "OpStatus" NOT NULL, -- Received, Accepted, Rejected
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_chekout"(
    "opid" "OperationId" NOT NULL,
    "to" "VersionId" NOT NULL
);

CREATE TABLE "op_obj_completion"(
    "opid" "OperationId" NOT NULL,
    "completed" boolean NOT NULL
);

CREATE TABLE "op_obj_content"(
    "opid" "OperationId" NOT NULL,
    "content" text
);

CREATE TABLE "op_obj_create"(
    "opid" "OperationId" NOT NULL,
    "content" text
);

CREATE TABLE "op_obj_reattach"(
    "opid" "OperationId" NOT NULL,
    "child" "ObjectiveId" NOT NULL,
    "newparent" "ObjectiveId" NOT NULL
);

CREATE TABLE "op_obj_reorder"(
    "opid" "OperationId" NOT NULL,
    "child" "ObjectiveId" NOT NULL,
    "moveafter" "ObjectiveId" NOT NULL
);

CREATE TABLE "op_transitive"(
    "opid" "OperationId" NOT NULL,
    "cause" "OperationId" NOT NULL
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

