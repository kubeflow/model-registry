CREATE TABLE IF NOT EXISTS "ExecutionProperty" (
  "execution_id" INTEGER NOT NULL,
  "name" TEXT NOT NULL,
  "is_custom_property" INTEGER NOT NULL,
  "int_value" INTEGER DEFAULT NULL,
  "double_value" REAL DEFAULT NULL,
  "string_value" TEXT,
  "byte_value" BLOB,
  "proto_value" BLOB,
  "bool_value" INTEGER DEFAULT NULL,
  PRIMARY KEY ("execution_id","name","is_custom_property")
);

CREATE INDEX "idx_execution_property_int" ON "ExecutionProperty" ("name","is_custom_property","int_value");
CREATE INDEX "idx_execution_property_double" ON "ExecutionProperty" ("name","is_custom_property","double_value");
CREATE INDEX "idx_execution_property_string" ON "ExecutionProperty" ("name","is_custom_property","string_value");