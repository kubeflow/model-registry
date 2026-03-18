-- Create Attribution table
CREATE TABLE IF NOT EXISTS "Attribution" (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    context_id INTEGER NOT NULL,
    artifact_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (context_id, artifact_id)
); 