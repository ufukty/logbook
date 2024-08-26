-- DROP DATABASE IF EXISTS logbook;
-- CREATE DATABASE logbook;
-- \c logbook;
-- ;
CREATE DOMAIN "AreaId" AS uuid;

CREATE DOMAIN "BookmarkId" AS uuid;

CREATE DOMAIN "CollaborationId" AS uuid;

CREATE DOMAIN "LinkId" AS uuid;

CREATE DOMAIN "ObjectiveId" AS uuid;

CREATE DOMAIN "OperationId" AS uuid;

CREATE DOMAIN "PropertiesId" AS uuid;

CREATE DOMAIN "UserId" AS uuid;

CREATE DOMAIN "VersionId" AS uuid;

CREATE DOMAIN "BottomUpPropsId" AS uuid;

CREATE TABLE "active"(
    "oid" "ObjectiveId" NOT NULL UNIQUE,
    "vid" "VersionId" NOT NULL
);

CREATE TABLE "objective"(
    "oid" "ObjectiveId" NOT NULL DEFAULT gen_random_uuid(),
    "vid" "VersionId" NOT NULL DEFAULT gen_random_uuid(),
    "based" "VersionId" NOT NULL, -- previous "vid" OR '00000000-0000-0000-0000-000000000000'
    "created_by" "OperationId" NOT NULL,
    "pid" "PropertiesId" NOT NULL,
    "bupid" "BottomUpPropsId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("oid", "vid")
);

CREATE INDEX "index_objective" ON "objective"("created_at");

CREATE TABLE "props"(
    "pid" "PropertiesId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "content" text NOT NULL,
    "completed" boolean NOT NULL,
    "creator" "UserId" NOT NULL,
    "owner" "UserId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
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
CREATE TABLE "bottom_up_props"(
    "bupid" "BottomUpPropsId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "children" int NOT NULL,
    "subtree_size" int NOT NULL,
    "subtree_completed" int NOT NULL
);

CREATE TABLE "bottom_up_props_third_person"(
    "bupid" "BottomUpPropsId" NOT NULL UNIQUE,
    "viewer" "UserId" NOT NULL,
    "children" int NOT NULL,
    "subtree_size" int NOT NULL,
    "subtree_completed" int NOT NULL
);

CREATE TYPE "OpType" AS ENUM(
    'checkout',
    'obj_completion',
    'obj_content',
    'obj_create_subtask',
    'obj_delete_subtask',
    'obj_attach',
    'obj_detach',
    'obj_reorder',
    'usr_register',
    'transitive',
    'double_transitive_merger'
);

CREATE TYPE "OpStatus" AS ENUM(
    'received',
    'accepted',
    'rejected'
);

CREATE TABLE "operation"(
    "opid" "OperationId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "subjectoid" "ObjectiveId" NOT NULL,
    "subjectvid" "VersionId" NOT NULL,
    "actor" "UserId" NOT NULL,
    "op_type" "OpType" NOT NULL, -- Transitive, Checkout, Completion, Reattach, Reorder, Content
    "op_status" "OpStatus" NOT NULL, -- Received, Accepted, Rejected
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_checkout"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "to" "VersionId" NOT NULL
);

CREATE TABLE "op_obj_completion"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "completed" boolean NOT NULL
);

CREATE TABLE "op_obj_content"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "content" text NOT NULL
);

CREATE TABLE "op_obj_create_subtask"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "content" text NOT NULL
);

CREATE TABLE "op_obj_delete_subtask"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "doid" "ObjectiveId" NOT NULL,
    "dvid" "VersionId" NOT NULL
);

CREATE TABLE "op_obj_attach"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "child" "ObjectiveId" NOT NULL,
    "newparent" "ObjectiveId" NOT NULL
);

CREATE TABLE "op_obj_detach"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "child" "ObjectiveId" NOT NULL,
    "newparent" "ObjectiveId" NOT NULL
);

CREATE TABLE "op_obj_reorder"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "child" "ObjectiveId" NOT NULL,
    "moveafter" "ObjectiveId" NOT NULL
);

CREATE TABLE "op_transitive"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "cause" "OperationId" NOT NULL
);

CREATE TABLE "op_double_transitive_merger"(
    "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "opid" "OperationId" NOT NULL,
    "first" "OperationId" NOT NULL,
    "second" "OperationId" NOT NULL
);

-- TODO: operations: reorder, note (create,update,delete), delegation (assign, unassign), collaboration (init, invite, restrict), versioning (rollback, fastforward)
;

CREATE TABLE "bookmark"(
    "bid" "BookmarkId" UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "oid" "ObjectiveId" NOT NULL,
    "title" text NOT NULL,
    "is_rock" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE "AreaType" AS ENUM(
    'solo',
    'collaboration'
);

CREATE TABLE "control_area"(
    "aid" "AreaId" UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    "root" "ObjectiveId" NOT NULL,
    "ar_type" "AreaType" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE TABLE "collaboration"(
    "cid" "CollaborationId" UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    "aid" "AreaId" NOT NULL,
    "creator" "UserId" NOT NULL,
    "admin" "UserId" NOT NULL,
    "leader" "UserId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE TABLE "collaborator"(
    "id" uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    "cid" "CollaborationId" NOT NULL,
    "uid" "UserId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

