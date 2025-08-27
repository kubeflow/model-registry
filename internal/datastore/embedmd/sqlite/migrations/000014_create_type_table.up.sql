CREATE TABLE IF NOT EXISTS "Type" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT NOT NULL,
  "version" TEXT DEFAULT NULL,
  "type_kind" INTEGER NOT NULL,
  "description" TEXT,
  "input_type" TEXT,
  "output_type" TEXT,
  "external_id" TEXT DEFAULT NULL,
  UNIQUE ("external_id")
);

CREATE INDEX "idx_type_name" ON "Type" ("name");
CREATE INDEX "idx_type_external_id" ON "Type" ("external_id");