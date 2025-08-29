-- Create Association table for SQLite
CREATE TABLE IF NOT EXISTS "Association" (
    id INTEGER PRIMARY KEY,
    context_id INTEGER NOT NULL,
    execution_id INTEGER NOT NULL,
    UNIQUE (context_id, execution_id)
);