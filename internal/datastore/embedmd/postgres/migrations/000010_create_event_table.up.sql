-- Create Event table
CREATE TABLE IF NOT EXISTS Event (
    id BIGSERIAL PRIMARY KEY,
    artifact_id BIGINT NOT NULL REFERENCES Artifact(id),
    execution_id BIGINT NOT NULL REFERENCES Execution(id),
    type VARCHAR(255) NOT NULL,
    milliseconds_since_epoch BIGINT,
    UNIQUE(artifact_id, execution_id, type)
); 