ARTIFACT_STATES = [
   "LIVE",
   "PENDING",
   "MARKED_FOR_DELETION",
   "DELETED",
   "ABANDONED",
   "REFERENCE",
   "UNKNOWN",
]
ARTIFACT_TYPE_PARAMS = [
    ("model-artifact", "s3://test-bucket/models/"),
    ("doc-artifact", "https://docs.example.com/docs/"),
    ("dataset-artifact", "s3://test-bucket/datasets/"),
    ("metric", "metrics://experiment/"),
    ("parameter", "params://experiment/"),
]
DEFAULT_API_TIMEOUT = 5.0

