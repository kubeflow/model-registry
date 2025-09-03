-- Add missing types for experiment tracking and metric history
INSERT INTO "Type" (name, version, type_kind, description, input_type, output_type, external_id)
SELECT t.* FROM (
    SELECT 'kf.MetricHistory' as name, NULL as version, 1 as type_kind, NULL as description, NULL as input_type, NULL as output_type, NULL as external_id
    UNION ALL SELECT 'kf.Experiment', NULL, 2, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'kf.ExperimentRun', NULL, 2, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'kf.DataSet', NULL, 1, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'kf.Metric', NULL, 1, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'kf.Parameter', NULL, 1, NULL, NULL, NULL, NULL
) t
WHERE NOT EXISTS (
    SELECT 1 FROM "Type"
    WHERE name = t.name AND version IS NULL
); 