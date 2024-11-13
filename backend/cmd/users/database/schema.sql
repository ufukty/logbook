CREATE DOMAIN "UserId" AS uuid;

CREATE TABLE "user"(
    "uid" "UserId" NOT NULL DEFAULT gen_random_uuid(),
    "deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

