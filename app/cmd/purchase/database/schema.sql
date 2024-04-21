CREATE DOMAIN ProductId AS uuid;

CREATE TABLE CREATE TABLE "purchase"(
    "purchaser" "UserId" NOT NULL,
    "product" "ProductId" NOT NULL,
    "authorization" text NOT NULL,
    "case" text NOT NULL,
    "completed_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
);

