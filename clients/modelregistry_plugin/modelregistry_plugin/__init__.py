"""
Model Registry MLflow Plugin

A MLflow tracking plugin that integrates with Kubeflow Model Registry.
"""

from .store import ModelRegistryStore

__version__ = "0.1.0"
__all__ = ["ModelRegistryStore"]
