"""Model registry storage backends.
"""

from .base import ProtoType, ScalarType
from .wrapper import MLMDStore

__all__ = [
    "ProtoType", "ScalarType",
    "MLMDStore",
]
