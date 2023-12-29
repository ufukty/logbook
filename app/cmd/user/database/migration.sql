DROP DATABASE IF EXISTS logbook_user;

CREATE DATABASE logbook_user;

CONNECT logbook_user;

CREATE TABLE
    "ACCESS" (
        "document_id" UUID NOT NULL REFERENCES "DOCUMENT" ("document_id"),
        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        "user-agent" VARCHAR(256),
        "ip-address" INET NOT NULL
    );
