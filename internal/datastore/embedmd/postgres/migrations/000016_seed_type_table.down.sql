DELETE FROM "Type" WHERE name IN (
    'mlmd.Dataset',
    'mlmd.Model',
    'mlmd.Metrics',
    'mlmd.Statistics',
    'mlmd.Train',
    'mlmd.Transform',
    'mlmd.Process',
    'mlmd.Evaluate',
    'mlmd.Deploy'
);
