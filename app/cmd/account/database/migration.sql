DROP DATABASE IF EXISTS logbook_user;

CREATE DATABASE logbook_user;

CONNECT logbook_user;

CREATE TABLE
    "user" (
        "uid" UUID NOT NULL DEFAULT gen_random_uuid(),
        "creation" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    "ACCESS" (
        "uid" UUID NOT NULL,
        "useragent" VARCHAR(256),
        "ipaddress" INET NOT NULL,
        "creation" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );