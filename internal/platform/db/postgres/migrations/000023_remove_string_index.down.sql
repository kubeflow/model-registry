CREATE INDEX IF NOT EXISTS idx_artifact_property_string ON "ArtifactProperty" (name, is_custom_property, string_value);
CREATE INDEX IF NOT EXISTS idx_context_property_string ON "ContextProperty" (name, is_custom_property, string_value);
CREATE INDEX IF NOT EXISTS idx_execution_property_string ON "ExecutionProperty" (name, is_custom_property, string_value);
