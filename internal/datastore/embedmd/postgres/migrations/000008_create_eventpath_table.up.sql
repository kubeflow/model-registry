-- Create EventPath table
CREATE TABLE IF NOT EXISTS "EventPath" (
    event_id INTEGER NOT NULL,
    is_index_step BOOLEAN NOT NULL,
    step_index INTEGER,
    step_key TEXT
);

CREATE INDEX idx_event_path_event_id ON "EventPath" (event_id); 