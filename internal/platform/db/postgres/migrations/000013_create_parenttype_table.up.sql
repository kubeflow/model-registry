-- Create ParentType table
CREATE TABLE IF NOT EXISTS "ParentType" (
    type_id INTEGER NOT NULL,
    parent_type_id INTEGER NOT NULL,
    PRIMARY KEY (type_id, parent_type_id)
); 