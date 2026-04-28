-- Create Artifact table
CREATE TABLE IF NOT EXISTS "Artifact" (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    type_id INTEGER NOT NULL,
    uri TEXT,
    state INTEGER DEFAULT NULL,
    name VARCHAR(255) DEFAULT NULL,
    external_id VARCHAR(255) DEFAULT NULL,
    create_time_since_epoch BIGINT NOT NULL DEFAULT '0',
    last_update_time_since_epoch BIGINT NOT NULL DEFAULT '0',
    PRIMARY KEY (id),
    UNIQUE (external_id),
    UNIQUE (type_id, name)
);

CREATE INDEX idx_artifact_uri ON "Artifact" (uri);
CREATE INDEX idx_artifact_create_time_since_epoch ON "Artifact" (create_time_since_epoch);
CREATE INDEX idx_artifact_last_update_time_since_epoch ON "Artifact" (last_update_time_since_epoch);
CREATE INDEX idx_artifact_external_id ON "Artifact" (external_id); 