-- Create ContextProperty table
CREATE TABLE IF NOT EXISTS ContextProperty (
    id BIGSERIAL PRIMARY KEY,
    context_id BIGINT NOT NULL REFERENCES Context(id),
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(255) NOT NULL,
    string_value TEXT,
    int_value BIGINT,
    double_value DOUBLE PRECISION,
    UNIQUE(context_id, name)
); 