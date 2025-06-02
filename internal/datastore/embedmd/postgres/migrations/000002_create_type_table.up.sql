-- Create Type table
CREATE TABLE IF NOT EXISTS Type (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255),
    external_id VARCHAR(255),
    description TEXT,
    UNIQUE(name, version)
); 