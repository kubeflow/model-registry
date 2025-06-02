-- Seed Type table
WITH types_to_insert (name, version, external_id, description) AS (
    VALUES
    ('mlmd.Dataset', NULL, NULL, NULL),
    ('mlmd.Model', NULL, NULL, NULL),
    ('mlmd.Metrics', NULL, NULL, NULL),
    ('mlmd.Statistics', NULL, NULL, NULL),
    ('mlmd.Train', NULL, NULL, NULL),
    ('mlmd.Transform', NULL, NULL, NULL),
    ('mlmd.Process', NULL, NULL, NULL),
    ('mlmd.Evaluate', NULL, NULL, NULL),
    ('mlmd.Deploy', NULL, NULL, NULL),
    ('kf.RegisteredModel', NULL, NULL, NULL),
    ('kf.ModelVersion', NULL, NULL, NULL),
    ('kf.DocArtifact', NULL, NULL, NULL),
    ('kf.ModelArtifact', NULL, NULL, NULL),
    ('kf.ServingEnvironment', NULL, NULL, NULL),
    ('kf.InferenceService', NULL, NULL, NULL),
    ('kf.ServeModel', NULL, NULL, NULL)
)
INSERT INTO type (name, version, external_id, description)
SELECT ti.name, ti.version, ti.external_id, ti.description
FROM types_to_insert ti
WHERE NOT EXISTS (
    SELECT 1 FROM type t_existing
    WHERE t_existing.name = ti.name AND (t_existing.version IS NULL AND ti.version IS NULL OR t_existing.version = ti.version)
)
ON CONFLICT (name, version) DO NOTHING; 