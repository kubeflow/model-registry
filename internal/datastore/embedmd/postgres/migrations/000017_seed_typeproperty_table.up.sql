-- Seed TypeProperty table
INSERT INTO typeproperty (type_id, name, data_type)
VALUES
    ((SELECT id FROM type WHERE name = 'kf.RegisteredModel'), 'description', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.RegisteredModel'), 'owner', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.RegisteredModel'), 'state', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelVersion'), 'author', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelVersion'), 'description', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelVersion'), 'model_name', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelVersion'), 'state', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelVersion'), 'version', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.DocArtifact'), 'description', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelArtifact'), 'description', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelArtifact'), 'model_format_name', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelArtifact'), 'model_format_version', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelArtifact'), 'service_account_name', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelArtifact'), 'storage_key', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ModelArtifact'), 'storage_path', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ServingEnvironment'), 'description', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.InferenceService'), 'description', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.InferenceService'), 'desired_state', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.InferenceService'), 'model_version_id', 'INT'),
    ((SELECT id FROM type WHERE name = 'kf.InferenceService'), 'registered_model_id', 'INT'),
    ((SELECT id FROM type WHERE name = 'kf.InferenceService'), 'runtime', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.InferenceService'), 'serving_environment_id', 'INT'),
    ((SELECT id FROM type WHERE name = 'kf.ServeModel'), 'description', 'STRING'),
    ((SELECT id FROM type WHERE name = 'kf.ServeModel'), 'model_version_id', 'INT')
ON CONFLICT (type_id, name) DO NOTHING; 