-- Create ArtifactProperty table
CREATE TABLE IF NOT EXISTS "ArtifactProperty" (
    artifact_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_custom_property BOOLEAN NOT NULL,
    int_value INTEGER DEFAULT NULL,
    double_value DOUBLE PRECISION DEFAULT NULL,
    string_value TEXT,
    byte_value BYTEA,
    proto_value BYTEA,
    bool_value BOOLEAN DEFAULT NULL,
    PRIMARY KEY (artifact_id, name, is_custom_property)
);

CREATE INDEX idx_artifact_property_int ON "ArtifactProperty" (name, is_custom_property, int_value);
CREATE INDEX idx_artifact_property_double ON "ArtifactProperty" (name, is_custom_property, double_value);
CREATE INDEX idx_artifact_property_string ON "ArtifactProperty" (name, is_custom_property, string_value); 