-- Create EventPath table
CREATE TABLE IF NOT EXISTS EventPath (
    event_id INTEGER NOT NULL,
    step_index INTEGER NOT NULL,
    is_index_step SMALLINT NOT NULL,
    execution_id INTEGER NOT NULL,
    PRIMARY KEY (step_index, execution_id)
);

CREATE INDEX idx_eventpath_event_id ON EventPath (event_id); 