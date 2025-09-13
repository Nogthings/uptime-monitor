-- +migrate Up
CREATE TABLE "status_checks" (
  "id" BIGSERIAL PRIMARY KEY,
  "service_id" BIGINT NOT NULL,
  "checked_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "status" VARCHAR(10) NOT NULL, -- 'up' or 'down'
  "status_code" INT,
  "response_time_ms" INT,
  "error_message" TEXT,
  CONSTRAINT fk_service
    FOREIGN KEY("service_id") 
    REFERENCES "services"("id")
    ON DELETE CASCADE
);

CREATE INDEX ON "status_checks" ("service_id", "checked_at" DESC);
