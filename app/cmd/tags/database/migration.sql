DROP DATABASE IF EXISTS logbook_tags;

CREATE DATABASE logbook_tags;

CONNECT logbook_tags;

CREATE TABLE
    "TAG" (
        "tid" UUID UNIQUE DEFAULT gen_random_UUID (), -- tag id
        "vid" -- version id of tag
        "text" TEXT NOT NULL,
        "uid" UUID NOT NULL,
        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    "TAGGING" (
        "tid" UUID NOT NULL,
        "oid" UUID NOT NULL, -- objective id
        "vid" UUID NOT NULL, -- 
        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    );
