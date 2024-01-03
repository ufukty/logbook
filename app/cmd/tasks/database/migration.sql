DROP DATABASE IF EXISTS logbook_objective;

CREATE DATABASE logbook_objective;

CONNECT logbook_objective;

CREATE TYPE OTYPE AS ENUM('rock', 'regular');

CREATE TABLE
    "objective" (
        "oid" UUID NOT NULL DEFAULT gen_random_uuid (), -- objective id
        "vid" UUID NOT NULL, -- version id
        "based" UUID, -- previous "vid" OR '00000000-0000-0000-0000-000000000000'
        "type" OTYPE NOT NULL DEFAULT 'regular',
        "content" TEXT NOT NULL,
        "creator" UUID NOT NULL, -- user id
        "creation" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY ("oid", "vid")
    );

CREATE INDEX "index_objective" ON "objective" ("creation");
    
CREATE TABLE
    "objective_link" (
        "lid" UUID NOT NULL DEFAULT generate_safe_uuid (), -- link id
        "sup_oid" UUID NOT NULL, -- super objective id
        "sup_vid" UUID NOT NULL, -- super version id
        "sub_oid" UUID NOT NULL, -- sub objective id
        "sub_vid" UUID NOT NULL, -- sub version id
        "creation" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY ("sup_oid", "sup_vid", "sub_oid")
    );

CREATE TABLE
    "objective_completion" (
        "oid" UUID NOT NULL,
        "vid" UUID NOT NULL,
        "actor" UUID NOT NULL, -- user id
        "completed" BOOLEAN NOT NULL DEFAULT FALSE,
        "creation" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    "objective_deleted" (
        "oid" UUID NOT NULL,
        "vid" UUID NOT NULL,
        "deletion" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    "computed_to_top" (
        "oid" UUID NOT NULL,
        "vid" UUID NOT NULL,
        "dependencies_are_cleared" BOOLEAN NOT NULL,
        "all_cleared" BOOLEAN NOT NULL,
        "degree" INT NOT NULL
        "completed_subtasks" INT NOT NULL,
        PRIMARY KEY ("oid", "vid")
    );

CREATE TABLE
    "computed_to_bottom" (
        "oid" UUID NOT NULL,
        "vid" UUID NOT NULL,
        "depth" INT NOT NULL,
        PRIMARY KEY ("oid", "vid")
    );

CREATE TABLE
    "objective_effective_version" (
        "oid" UUIT NOT NULL UNIQUE,
        "vid" UUIT NOT NULL
    )

CREATE INDEX "index_effective_version" ON "objective_effective_version" ("oid");