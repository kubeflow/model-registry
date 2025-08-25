import os
from urllib.parse import urlparse

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

REGISTRY_URL = os.environ.get("MR_URL", "http://localhost:8080")
parsed = urlparse(REGISTRY_URL)
host, port = parsed.netloc.split(":")
REGISTRY_HOST = f"{parsed.scheme}://{host}"
REGISTRY_PORT = int(port)

MAX_POLL_TIME = 10
POLL_INTERVAL = 1

