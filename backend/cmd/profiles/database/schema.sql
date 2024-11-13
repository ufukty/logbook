-- DROP DATABASE IF EXISTS logbook_profiles;
-- CREATE DATABASE logbook_profiles;
-- CONNECT logbook_profiles;
CREATE DOMAIN "UserId" AS uuid;

CREATE DOMAIN "UserId" AS uuid;

CREATE TABLE "user"(
    "uid" "UserId" NOT NULL DEFAULT gen_random_uuid(),
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE DOMAIN "HumanName" AS text;

CREATE TABLE "profile"(
    "uid" "UserId" NOT NULL,
    "firstname" "HumanName" NOT NULL,
    "lastname" "HumanName" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

