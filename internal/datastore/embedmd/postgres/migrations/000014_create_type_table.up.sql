-- Create Type table
CREATE TABLE IF NOT EXISTS "Type" (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255) DEFAULT NULL,
    type_kind SMALLINT NOT NULL,
    description TEXT,
    input_type TEXT,
    output_type TEXT,
    external_id VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE (external_id)
);

CREATE INDEX idx_type_name ON "Type" (name);
CREATE INDEX idx_type_external_id ON "Type" (external_id); 