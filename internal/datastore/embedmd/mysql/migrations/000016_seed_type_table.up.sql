INSERT INTO `Type` (`name`, `version`, `type_kind`, `description`, `input_type`, `output_type`, `external_id`)
SELECT t.* FROM (
    SELECT 'mlmd.Dataset' as name, NULL as version, 1 as type_kind, NULL as description, NULL as input_type, NULL as output_type, NULL as external_id
    UNION ALL SELECT 'mlmd.Model', NULL, 1, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'mlmd.Metrics', NULL, 1, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'mlmd.Statistics', NULL, 1, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'mlmd.Train', NULL, 0, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'mlmd.Transform', NULL, 0, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'mlmd.Process', NULL, 0, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'mlmd.Evaluate', NULL, 0, NULL, NULL, NULL, NULL
    UNION ALL SELECT 'mlmd.Deploy', NULL, 0, NULL, NULL, NULL, NULL
) t
WHERE NOT EXISTS (
    SELECT 1 FROM `Type`
    WHERE `name` = t.name AND `version` IS NULL
);
