CREATE MATERIALIZED VIEW IF NOT EXISTS context_property_options AS
    SELECT type_id, name, is_custom_property, array_agg(string_value) string_value, NULL::text[] array_value, NULL::float8 min_double_value, NULL::float8 max_double_value, NULL::integer min_int_value, NULL::integer max_int_value FROM (
            SELECT DISTINCT "Context".type_id, "ContextProperty".name, is_custom_property, "ContextProperty".string_value
            FROM "Context"
                INNER JOIN "ContextProperty" ON "Context".id="ContextProperty".context_id
            WHERE string_value IS NOT NULL AND
                    string_value != '' AND
                    string_value IS NOT JSON ARRAY
            ORDER BY string_value
    )
    GROUP BY type_id, name, is_custom_property HAVING MAX(CHAR_LENGTH(string_value)) <= 100 AND COUNT(*) < 50

    UNION

    SELECT type_id, name, is_custom_property, NULL, array_agg(string_value), NULL::float8 min_double_value, NULL::float8 max_double_value, NULL::integer min_int_value, NULL::integer max_int_value FROM (
            SELECT DISTINCT "Context".type_id, "ContextProperty".name, is_custom_property, json_array_elements_text("ContextProperty".string_value::json) AS string_value
            FROM "Context"
                INNER JOIN "ContextProperty" ON "Context".id="ContextProperty".context_id
            WHERE string_value IS JSON ARRAY
    )
    GROUP BY type_id, name, is_custom_property HAVING MAX(CHAR_LENGTH(string_value)) <= 100 AND COUNT(*) < 50

    UNION

    SELECT type_id, name, is_custom_property, NULL, NULL, min(double_value), max(double_value), NULL, NULL FROM (
        SELECT
            "Context".type_id, "ContextProperty".name, "ContextProperty".is_custom_property, "ContextProperty".double_value
        FROM "Context"
            INNER JOIN "ContextProperty" ON "Context".id="ContextProperty".context_id
        WHERE double_value IS NOT NULL
    ) values GROUP BY type_id, name, is_custom_property

    UNION

    SELECT type_id, name, is_custom_property, NULL, NULL, NULL, NULL, min(int_value), max(int_value) FROM (
        SELECT "Context".type_id, "ContextProperty".name, is_custom_property, "ContextProperty".int_value
        FROM "Context"
            INNER JOIN "ContextProperty" ON "Context".id="ContextProperty".context_id
        WHERE int_value IS NOT NULL
    ) values GROUP BY type_id, name, is_custom_property
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS artifact_property_options AS
    SELECT type_id, name, is_custom_property, array_agg(string_value) string_value, NULL::text[] array_value, NULL::float8 min_double_value, NULL::float8 max_double_value, NULL::integer min_int_value, NULL::integer max_int_value FROM (
            SELECT DISTINCT "Artifact".type_id, "ArtifactProperty".name, is_custom_property, "ArtifactProperty".string_value
            FROM "Artifact"
                INNER JOIN "ArtifactProperty" ON "Artifact".id="ArtifactProperty".artifact_id
            WHERE string_value IS NOT NULL AND
                    string_value != '' AND
                    string_value IS NOT JSON ARRAY
            ORDER BY string_value
    )
    GROUP BY type_id, name, is_custom_property HAVING MAX(CHAR_LENGTH(string_value)) <= 100 AND COUNT(*) < 50

    UNION

    SELECT type_id, name, is_custom_property, NULL, array_agg(string_value), NULL::float8 min_double_value, NULL::float8 max_double_value, NULL::integer min_int_value, NULL::integer max_int_value FROM (
            SELECT DISTINCT "Artifact".type_id, "ArtifactProperty".name, is_custom_property, json_array_elements_text("ArtifactProperty".string_value::json) AS string_value
            FROM "Artifact"
                INNER JOIN "ArtifactProperty" ON "Artifact".id="ArtifactProperty".artifact_id
            WHERE string_value IS JSON ARRAY
    )
    GROUP BY type_id, name, is_custom_property HAVING MAX(CHAR_LENGTH(string_value)) <= 100 AND COUNT(*) < 50

    UNION

    SELECT type_id, name, is_custom_property, NULL, NULL, min(double_value), max(double_value), NULL, NULL FROM (
        SELECT
            "Artifact".type_id, "ArtifactProperty".name, "ArtifactProperty".is_custom_property, "ArtifactProperty".double_value
        FROM "Artifact"
            INNER JOIN "ArtifactProperty" ON "Artifact".id="ArtifactProperty".artifact_id
        WHERE double_value IS NOT NULL
    ) values GROUP BY type_id, name, is_custom_property

    UNION

    SELECT type_id, name, is_custom_property, NULL, NULL, NULL, NULL, min(int_value), max(int_value) FROM (
        SELECT "Artifact".type_id, "ArtifactProperty".name, is_custom_property, "ArtifactProperty".int_value
        FROM "Artifact"
            INNER JOIN "ArtifactProperty" ON "Artifact".id="ArtifactProperty".artifact_id
        WHERE int_value IS NOT NULL
    ) values GROUP BY type_id, name, is_custom_property
WITH NO DATA;
