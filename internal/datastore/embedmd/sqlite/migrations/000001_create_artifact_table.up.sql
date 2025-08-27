-- Create Artifact table for SQLite
CREATE TABLE IF NOT EXISTS "Artifact" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id INTEGER NOT NULL,
    uri TEXT,
    state INTEGER DEFAULT NULL,
    name TEXT DEFAULT NULL,
    external_id TEXT DEFAULT NULL,
    create_time_since_epoch INTEGER NOT NULL DEFAULT '0',
    last_update_time_since_epoch INTEGER NOT NULL DEFAULT '0',
    UNIQUE (external_id),
    UNIQUE (type_id, name)
);

CREATE INDEX idx_artifact_uri ON "Artifact" (uri);
CREATE INDEX idx_artifact_create_time_since_epoch ON "Artifact" (create_time_since_epoch);
CREATE INDEX idx_artifact_last_update_time_since_epoch ON "Artifact" (last_update_time_since_epoch);
CREATE INDEX idx_artifact_external_id ON "Artifact" (external_id);