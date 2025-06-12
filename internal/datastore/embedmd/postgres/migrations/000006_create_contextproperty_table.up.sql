-- Create ContextProperty table
CREATE TABLE IF NOT EXISTS "ContextProperty" (
    context_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_custom_property BOOLEAN NOT NULL,
    int_value INTEGER DEFAULT NULL,
    double_value DOUBLE PRECISION DEFAULT NULL,
    string_value TEXT,
    byte_value BYTEA,
    proto_value BYTEA,
    bool_value BOOLEAN DEFAULT NULL,
    PRIMARY KEY (context_id, name, is_custom_property)
);

CREATE INDEX idx_context_property_int ON "ContextProperty" (name, is_custom_property, int_value);
CREATE INDEX idx_context_property_double ON "ContextProperty" (name, is_custom_property, double_value);
CREATE INDEX idx_context_property_string ON "ContextProperty" (name, is_custom_property, string_value); 