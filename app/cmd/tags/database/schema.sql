-- DROP DATABASE IF EXISTS logbook_tags;
-- CREATE DATABASE logbook_tags;
-- CONNECT logbook_tags;
CREATE DOMAIN "ObjectiveId" AS uuid;

CREATE DOMAIN "TagId" AS uuid;

CREATE DOMAIN "UserId" AS uuid;

CREATE DOMAIN "VersionId" AS uuid;

CREATE TABLE "tag"(
  "tid" "TagId" UNIQUE DEFAULT gen_random_UUID(), -- tag id
  "vid" "VersionId" NOT NULL, -- version id of tag
  "text" text NOT NULL,
  "uid" "UserId" NOT NULL,
  "deleted" boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "tagging"(
  "tid" "TagId" NOT NULL,
  "oid" "ObjectiveId" NOT NULL,
  "vid" "VersionId" NOT NULL,
  "deleted" boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

