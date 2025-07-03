-- Create ParentContext table
CREATE TABLE IF NOT EXISTS "ParentContext" (
    context_id INTEGER NOT NULL,
    parent_context_id INTEGER NOT NULL,
    PRIMARY KEY (context_id, parent_context_id)
);

CREATE INDEX idx_parentcontext_parent_context_id ON "ParentContext" (parent_context_id); 