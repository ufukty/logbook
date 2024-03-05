-- DROP DATABASE IF EXISTS logbook_user;
-- CREATE DATABASE logbook_user;
-- CONNECT logbook_user;
CREATE TABLE "user"(
    "uid" uuid NOT NULL DEFAULT gen_random_uuid(),
    "email" text NOT NULL,
    "hash" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "access"(
    "uid" uuid NOT NULL,
    "useragent" text,
    "ipaddress" inet NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

