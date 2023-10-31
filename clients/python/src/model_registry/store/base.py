from typing import Union

from ml_metadata.proto import Artifact, Context


ProtoType = Union[Artifact, Context]
"""Union of all proto types. """
ScalarType = Union[str, int, float, bool]
"""Union of all scalar types.

Those types are easy to map to and from proto types, and can also be queried in MLMD.
"""
