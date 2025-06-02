-- Create TypeProperty table
CREATE TABLE IF NOT EXISTS TypeProperty (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL REFERENCES Type(id),
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(255) NOT NULL,
    description TEXT,
    UNIQUE(type_id, name)
); 