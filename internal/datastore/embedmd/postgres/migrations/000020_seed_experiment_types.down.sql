DELETE FROM "Type" WHERE name IN (
    'kf.MetricHistory',
    'kf.Experiment',
    'kf.ExperimentRun',
    'kf.DataSet',
    'kf.Metric',
    'kf.Parameter'
); 