-- Migrate existing catalog model names to namespaced format (source_id:model_name).
-- Catalog models are Context rows with type_id = kf.CatalogModel; source_id is in ContextProperty.
-- Only update names that do not already contain ':' (legacy format).
UPDATE "Context" c
SET name = cp.string_value || ':' || c.name
FROM "ContextProperty" cp
WHERE cp.context_id = c.id
  AND cp.name = 'source_id'
  AND cp.is_custom_property = false
  AND c.type_id = (SELECT id FROM "Type" WHERE name = 'kf.CatalogModel' LIMIT 1)
  AND c.name NOT LIKE '%:%';
