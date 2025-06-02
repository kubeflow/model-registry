-- Create ExecutionProperty table
CREATE TABLE IF NOT EXISTS ExecutionProperty (
    id BIGSERIAL PRIMARY KEY,
    execution_id BIGINT NOT NULL REFERENCES Execution(id),
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(255) NOT NULL,
    string_value TEXT,
    int_value BIGINT,
    double_value DOUBLE PRECISION,
    UNIQUE(execution_id, name)
); 