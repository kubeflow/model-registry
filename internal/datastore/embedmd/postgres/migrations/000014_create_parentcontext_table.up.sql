-- Create ParentContext table
CREATE TABLE IF NOT EXISTS ParentContext (
    id BIGSERIAL PRIMARY KEY,
    context_id BIGINT NOT NULL REFERENCES Context(id),
    parent_context_id BIGINT NOT NULL REFERENCES Context(id),
    UNIQUE(context_id, parent_context_id)
); 