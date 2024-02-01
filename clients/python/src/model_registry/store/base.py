"""Base classes and types for MLMD store."""

from typing import Union

from ml_metadata.proto import Artifact, Context

# Union of all proto types.
ProtoType = Union[Artifact, Context]
# Union of all scalar types.
#
# Those types are easy to map to and from proto types, and can also be queried in MLMD.
ScalarType = Union[str, int, float, bool]
