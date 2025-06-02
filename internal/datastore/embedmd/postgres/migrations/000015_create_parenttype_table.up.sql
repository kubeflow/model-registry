-- Create ParentType table
CREATE TABLE IF NOT EXISTS ParentType (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL REFERENCES Type(id),
    parent_type_id BIGINT NOT NULL REFERENCES Type(id),
    UNIQUE(type_id, parent_type_id)
); 