-- Create Context table
CREATE TABLE IF NOT EXISTS "Context" (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    type_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    external_id VARCHAR(255) DEFAULT NULL,
    create_time_since_epoch BIGINT NOT NULL DEFAULT '0',
    last_update_time_since_epoch BIGINT NOT NULL DEFAULT '0',
    PRIMARY KEY (id),
    UNIQUE (type_id, name),
    UNIQUE (external_id)
);

CREATE INDEX idx_context_create_time_since_epoch ON "Context" (create_time_since_epoch);
CREATE INDEX idx_context_last_update_time_since_epoch ON "Context" (last_update_time_since_epoch);
CREATE INDEX idx_context_external_id ON "Context" (external_id); 