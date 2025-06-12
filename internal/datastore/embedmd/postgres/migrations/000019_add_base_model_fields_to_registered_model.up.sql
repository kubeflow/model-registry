INSERT INTO "TypeProperty" (type_id, name, data_type)
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'language', 4 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'library_name', 3 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'license_link', 3 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'license', 3 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'logo', 3 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'maturity', 3 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'provider', 3 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'readme', 3 UNION ALL
SELECT (SELECT id FROM "Type" WHERE name = 'kf.RegisteredModel'), 'tasks', 4; 