-- Create Association table
CREATE TABLE IF NOT EXISTS "Association" (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    context_id INTEGER NOT NULL,
    execution_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (context_id, execution_id)
); 