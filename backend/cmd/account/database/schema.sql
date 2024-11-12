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

CREATE DOMAIN "Email" AS text;

CREATE TABLE "login"(
    "lid" "LoginId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "email" "Email" NOT NULL,
    "hash" text NOT NULL,
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

CREATE DOMAIN "AccessId" AS uuid;

CREATE DOMAIN "UserAgent" AS text;

CREATE TABLE "access"(
    "aid" "AccessId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "useragent" "UserAgent" NOT NULL,
    "ipaddress" inet NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE DOMAIN "SessionId" AS uuid;

-- base64 encoded
CREATE DOMAIN "SessionToken" AS text;

CREATE TABLE "session_standard"(
    "sid" "SessionId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "token" "SessionToken" NOT NULL,
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "session_account_read"(
    "sid" "SessionId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "token" "SessionToken" NOT NULL,
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "session_account_write"(
    "sid" "SessionId" NOT NULL DEFAULT gen_random_uuid(),
    "uid" "UserId" NOT NULL,
    "token" "SessionToken" NOT NULL,
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

