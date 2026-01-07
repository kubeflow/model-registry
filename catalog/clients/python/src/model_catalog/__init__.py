"""Main package for the Kubeflow model catalog client."""

__version__ = "0.1.0"

from ._client import (
    CatalogAPIClient,
    CatalogAPIError,
    CatalogConnectionError,
    CatalogError,
    CatalogNotFoundError,
    CatalogValidationError,
)

__all__ = [
    "CatalogAPIClient",
    "CatalogAPIError",
    "CatalogConnectionError",
    "CatalogError",
    "CatalogNotFoundError",
    "CatalogValidationError",
]
