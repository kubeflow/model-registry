"""Model Registry MLflow Plugin.

A MLflow tracking plugin that integrates with Kubeflow Model Registry.
"""

from .tracking_store import ModelRegistryTrackingStore

__version__ = "0.1.0"
__all__ = ["ModelRegistryTrackingStore"]
