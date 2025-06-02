DELETE FROM type WHERE name IN (
    'mlmd.Dataset',
    'mlmd.Model',
    'mlmd.Metrics',
    'mlmd.Statistics',
    'mlmd.Train',
    'mlmd.Transform',
    'mlmd.Process',
    'mlmd.Evaluate',
    'mlmd.Deploy',
    'kf.RegisteredModel',
    'kf.ModelVersion',
    'kf.DocArtifact',
    'kf.ModelArtifact',
    'kf.ServingEnvironment',
    'kf.InferenceService',
    'kf.ServeModel'
); 