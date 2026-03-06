-- Revert catalog model names from namespaced format (source_id:model_name) back to model_name only.
-- Only reverts rows where the prefix matches the context's source_id property.
UPDATE "Context" c
SET name = substring(c.name from position(':' in c.name) + 1)
FROM "ContextProperty" cp
WHERE cp.context_id = c.id
  AND cp.name = 'source_id'
  AND cp.is_custom_property = false
  AND c.type_id = (SELECT id FROM "Type" WHERE name = 'kf.CatalogModel' LIMIT 1)
  AND c.name LIKE cp.string_value || ':%';
