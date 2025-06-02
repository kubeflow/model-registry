-- Create Attribution table
CREATE TABLE IF NOT EXISTS Attribution (
    id BIGSERIAL PRIMARY KEY,
    context_id BIGINT NOT NULL REFERENCES Context(id),
    artifact_id BIGINT NOT NULL REFERENCES Artifact(id),
    UNIQUE(context_id, artifact_id)
); 