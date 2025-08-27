CREATE TABLE IF NOT EXISTS "Context" (
  "id" INTEGER PRIMARY KEY,
  "type_id" INTEGER NOT NULL,
  "name" TEXT NOT NULL,
  "external_id" TEXT DEFAULT NULL,
  "create_time_since_epoch" BIGINT NOT NULL DEFAULT 0,
  "last_update_time_since_epoch" BIGINT NOT NULL DEFAULT 0,
  UNIQUE ("type_id","name"),
  UNIQUE ("external_id")
);

CREATE INDEX "idx_context_create_time_since_epoch" ON "Context" ("create_time_since_epoch");
CREATE INDEX "idx_context_last_update_time_since_epoch" ON "Context" ("last_update_time_since_epoch");
CREATE INDEX "idx_context_external_id" ON "Context" ("external_id");