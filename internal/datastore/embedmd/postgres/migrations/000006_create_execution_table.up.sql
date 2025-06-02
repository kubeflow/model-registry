-- Create Execution table
CREATE TABLE IF NOT EXISTS Execution (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL REFERENCES Type(id),
    last_known_state VARCHAR(255),
    name VARCHAR(255),
    external_id VARCHAR(255),
    create_time_since_epoch BIGINT,
    last_update_time_since_epoch BIGINT,
    UNIQUE(external_id)
); 