-- +migrate Up
CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "services" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "target" VARCHAR(255) NOT NULL,
  "check_interval_seconds" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  CONSTRAINT fk_user
    FOREIGN KEY("user_id") 
    REFERENCES "users"("id")
    ON DELETE CASCADE
);
