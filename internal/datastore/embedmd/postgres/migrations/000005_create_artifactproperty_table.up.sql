-- Create ArtifactProperty table
CREATE TABLE IF NOT EXISTS ArtifactProperty (
    id BIGSERIAL PRIMARY KEY,
    artifact_id BIGINT NOT NULL REFERENCES Artifact(id),
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(255) NOT NULL,
    string_value TEXT,
    int_value BIGINT,
    double_value DOUBLE PRECISION,
    UNIQUE(artifact_id, name)
); 