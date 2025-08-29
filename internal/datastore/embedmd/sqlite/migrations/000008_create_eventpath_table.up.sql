CREATE TABLE IF NOT EXISTS "EventPath" (
  "event_id" INTEGER NOT NULL,
  "is_index_step" BOOLEAN NOT NULL,
  "step_index" INTEGER DEFAULT NULL,
  "step_key" TEXT
);

CREATE INDEX "idx_eventpath_event_id" ON "EventPath" ("event_id");