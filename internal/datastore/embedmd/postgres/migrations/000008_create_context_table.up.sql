-- Create Context table
CREATE TABLE IF NOT EXISTS Context (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL REFERENCES Type(id),
    name VARCHAR(255) NOT NULL,
    external_id VARCHAR(255),
    create_time_since_epoch BIGINT,
    last_update_time_since_epoch BIGINT,
    UNIQUE(external_id)
); 