-- Create Event table
CREATE TABLE IF NOT EXISTS "Event" (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    artifact_id INTEGER NOT NULL,
    execution_id INTEGER NOT NULL,
    type INTEGER NOT NULL,
    milliseconds_since_epoch BIGINT DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE (artifact_id, execution_id, type)
);

CREATE INDEX idx_event_execution_id ON "Event" (execution_id); 