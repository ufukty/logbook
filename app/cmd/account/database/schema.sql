-- DROP DATABASE IF EXISTS logbook_user;
-- CREATE DATABASE logbook_user;
-- CONNECT logbook_user;
CREATE DOMAIN "UserId" AS uuid;

CREATE TABLE "user"(
    "uid" "UserId" NOT NULL DEFAULT gen_random_uuid(),
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE DOMAIN "LoginId" AS uuid;

CREATE TABLE "login"(
    "lid" "LoginId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "email" text NOT NULL,
    "hash" text NOT NULL,
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "profile"(
    "uid" "UserId" NOT NULL,
    "firstname" text NOT NULL,
    "lastname" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE DOMAIN "AccessId" AS uuid;

CREATE TABLE "access"(
    "aid" "AccessId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "useragent" text,
    "ipaddress" inet NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE DOMAIN "SessionId" AS uuid;

CREATE TABLE "session_standard"(
    "sid" "SessionId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "token" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "session_account_read"(
    "sid" "SessionId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "token" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "session_account_write"(
    "sid" "SessionId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "token" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);