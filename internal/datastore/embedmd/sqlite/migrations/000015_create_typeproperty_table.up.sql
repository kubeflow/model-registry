CREATE TABLE IF NOT EXISTS "TypeProperty" (
  "type_id" INTEGER NOT NULL,
  "name" TEXT NOT NULL,
  "data_type" INTEGER DEFAULT NULL,
  PRIMARY KEY ("type_id","name")
);