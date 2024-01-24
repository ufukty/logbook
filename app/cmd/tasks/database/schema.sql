-- DROP DATABASE IF EXISTS logbook_objective;
-- CREATE DATABASE logbook_objective;
-- \c logbook_objective;
--
CREATE TABLE "versioning_config" (
    "oid" uuid NOT NULL, -- objective id
    "first" uuid NOT NULL, -- version id
    "effective" uuid NOT NULL -- version id
);

CREATE TABLE "version" (
    "vid" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid (), -- version id
    "based" uuid -- version id
);

CREATE TYPE OTYPE AS ENUM (
    'rock',
    'regular'
);

CREATE TABLE "objective" (
    "oid" uuid NOT NULL DEFAULT gen_random_uuid (), -- objective id
    "vid" uuid NOT NULL, -- version id
    "based" uuid, -- previous "vid" OR '00000000-0000-0000-0000-000000000000'
    "type" OTYPE NOT NULL DEFAULT 'regular',
    "content" text NOT NULL,
    "creator" uuid NOT NULL, -- user id
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("oid", "vid")
);

CREATE INDEX "index_objective" ON "objective" ("created_at");

CREATE TABLE "objective_link" (
    "lid" uuid NOT NULL DEFAULT gen_random_uuid (), -- link id
    "sup_oid" uuid NOT NULL, -- super objective id
    "sup_vid" uuid NOT NULL, -- super version id
    "sub_oid" uuid NOT NULL, -- sub objective id
    "sub_vid" uuid NOT NULL, -- sub version id
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("sup_oid", "sup_vid", "sub_oid")
);

CREATE TABLE "objective_completion" (
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL,
    "actor" uuid NOT NULL, -- user id
    "completed" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "objective_deleted" (
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL,
    "deletion" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Computed properties and user preferences per item per user
CREATE TABLE "objective_view" (
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL,
    "uid" uuid NOT NULL, -- viewer
    "degree" int NOT NULL,
    "depth" int NOT NULL,
    "ready" boolean NOT NULL,
    "completion_pct" float NOT NULL,
    "fold" boolean NOT NULL
);

CREATE TABLE "computed_to_top" (
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL,
    "dependencies_are_cleared" boolean NOT NULL,
    "all_cleared" boolean NOT NULL,
    "degree" int NOT NULL,
    "completed_subtasks" int NOT NULL,
    PRIMARY KEY ("oid", "vid")
);

CREATE TABLE "computed_to_bottom" (
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL,
    "depth" int NOT NULL,
    PRIMARY KEY ("oid", "vid")
);

CREATE TABLE "objective_effective_version" (
    "oid" uuid NOT NULL UNIQUE,
    "vid" uuid NOT NULL
);

CREATE INDEX "index_effective_version" ON "objective_effective_version" ("oid");

CREATE TABLE "op_objective_create" (
    "opid" uuid NOT NULL DEFAULT gen_random_uuid (),
    "poid" uuid NOT NULL, -- parent oid
    "pvid" uuid NOT NULL, -- parent vid based on
    "actor" uuid NOT NULL,
    "content" text,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_delete" (
    "opid" uuid NOT NULL DEFAULT gen_random_uuid (),
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL, -- based on
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_content_update" (
    "opid" uuid NOT NULL DEFAULT gen_random_uuid (),
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL, -- based on
    "content" text,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_attach_subobjective" (
    "opid" uuid NOT NULL DEFAULT gen_random_uuid (),
    "sup_oid" uuid NOT NULL,
    "sup_vid" uuid NOT NULL, -- based on
    "sub_oid" uuid NOT NULL,
    "sub_vid" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "op_objective_update_completion" (
    "opid" uuid NOT NULL DEFAULT gen_random_uuid (),
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL,
    "completed" boolean NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- TODO: operations: reorder, delegation (assign, unassign), collaboration (init, invite, restrict), versioning (rollback, fastforward)
;

CREATE TABLE "bookmark" (
    "user" uuid NOT NULL,
    "oid" uuid NOT NULL,
    "vid" uuid NOT NULL,
    "name" text,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

