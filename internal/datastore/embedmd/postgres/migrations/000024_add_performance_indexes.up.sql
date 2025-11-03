-- Add performance indexes for Context and Attribution tables
-- These indexes optimize queries that filter by type_id and join through Attribution

-- Index on Context.type_id for efficient filtering
CREATE INDEX IF NOT EXISTS idx_context_type_id ON "Context" (type_id);

-- Composite index on Attribution for efficient joins with Context and ArtifactProperty
CREATE INDEX IF NOT EXISTS idx_attribution_context_artifact ON "Attribution" (context_id, artifact_id);

-- Index on ArtifactProperty.artifact_id for efficient joins with Attribution
CREATE INDEX IF NOT EXISTS idx_artifact_property_artifact_id ON "ArtifactProperty" (artifact_id);