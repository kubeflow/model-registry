DELETE FROM `TypeProperty` WHERE type_id=(
    SELECT id FROM `Type` WHERE name = 'kf.RegisteredModel'
) AND `name` IN (
    'language',
    'library_name',
    'license_link',
    'license',
    'logo',
    'maturity',
    'provider',
    'readme',
    'tasks'
);