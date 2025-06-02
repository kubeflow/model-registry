-- Create Attribution table
CREATE TABLE IF NOT EXISTS Attribution (
    id SERIAL,
    context_id INTEGER NOT NULL,
    artifact_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (context_id, artifact_id)
); 