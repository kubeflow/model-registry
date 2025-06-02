-- Create EventPath table
CREATE TABLE IF NOT EXISTS EventPath (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES Event(id),
    is_index_step BOOLEAN NOT NULL,
    step_index INTEGER,
    step_key VARCHAR(255)
); 