-- Create Association table
CREATE TABLE IF NOT EXISTS Association (
    id BIGSERIAL PRIMARY KEY,
    context_id BIGINT NOT NULL REFERENCES Context(id),
    execution_id BIGINT NOT NULL REFERENCES Execution(id),
    UNIQUE(context_id, execution_id)
); 