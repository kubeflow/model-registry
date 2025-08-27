CREATE TABLE IF NOT EXISTS "Attribution" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "context_id" INTEGER NOT NULL,
  "artifact_id" INTEGER NOT NULL,
  UNIQUE ("context_id","artifact_id")
);