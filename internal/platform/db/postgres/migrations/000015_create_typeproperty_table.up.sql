-- Create TypeProperty table
CREATE TABLE IF NOT EXISTS "TypeProperty" (
    type_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_type INTEGER DEFAULT NULL,
    PRIMARY KEY (type_id, name)
); 