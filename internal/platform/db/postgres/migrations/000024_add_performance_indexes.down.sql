-- Remove performance indexes added in 000024_add_performance_indexes.up.sql

DROP INDEX IF EXISTS idx_context_type_id;
DROP INDEX IF EXISTS idx_attribution_context_artifact;
DROP INDEX IF EXISTS idx_artifact_property_artifact_id;