-- Create Association table
CREATE TABLE IF NOT EXISTS Association (
    id SERIAL,
    context_id INTEGER NOT NULL,
    execution_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (context_id, execution_id)
); 