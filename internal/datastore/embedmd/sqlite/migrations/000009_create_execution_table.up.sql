CREATE TABLE IF NOT EXISTS "Execution" (
  "id" INTEGER PRIMARY KEY,
  "type_id" INTEGER NOT NULL,
  "last_known_state" INTEGER DEFAULT NULL,
  "name" TEXT DEFAULT NULL,
  "external_id" TEXT DEFAULT NULL,
  "create_time_since_epoch" BIGINT NOT NULL DEFAULT 0,
  "last_update_time_since_epoch" BIGINT NOT NULL DEFAULT 0,
  UNIQUE ("external_id"),
  UNIQUE ("type_id","name")
);

CREATE INDEX "idx_execution_create_time_since_epoch" ON "Execution" ("create_time_since_epoch");
CREATE INDEX "idx_execution_last_update_time_since_epoch" ON "Execution" ("last_update_time_since_epoch");
CREATE INDEX "idx_execution_external_id" ON "Execution" ("external_id");