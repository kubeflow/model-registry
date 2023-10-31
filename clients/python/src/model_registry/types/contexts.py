"""Context types for model registry.

Contexts group related Artifacts together.
They provide a way to organize and categorize components in a workflow.

Those types are used to map between proto types based on contexts and Python objects.

TODO:
    * Move part of the description to API Reference docs (#120).
"""

from abc import ABC
from typing import Union

from .artifacts import BaseArtifact, ModelArtifact
from .base import ProtoBase

from attrs import define, field


ScalarType = Union[str, int, float, bool]


@define(slots=False, init=False)
class BaseContext(ProtoBase, ABC):
    """Abstract base class for all contexts."""


@define(slots=False)
class ModelVersion(BaseContext):
    """Represents a model version.

    Attributes:
        model (ModelArtifact): Model associated with this version.
        version (str): Version of the model.
        author (str): Author of the model version.
        description (str, optional): Description of the object.
        external_id (str, optional): Customizable ID. Has to be unique among instances of the same type.
        artifacts (list[BaseArtifact]): Artifacts associated with this version.
        tags (list[str]): Tags associated with this version.
        metadata (dict[str, ScalarType]): Metadata associated with this version.
    """

    model: ModelArtifact
    version: str
    author: str
    artifacts: list[BaseArtifact] = field(init=False, factory=list)
    tags: list[str] = field(init=False, factory=list)
    metadata: dict[str, ScalarType] = field(init=False, factory=dict)

    def __attrs_post_init__(self) -> None:
        self.name = self.version


@define(slots=False)
class RegisteredModel(BaseContext):
    """Represents a registered model.

    Attributes:
        name (str): Registered model name.
        description (str, optional): Description of the object.
        external_id (str, optional): Customizable ID. Has to be unique among instances of the same type.
        versions (list[ModelVersion]): Versions associated with this model.
    """

    name: str
    versions: list[ModelVersion] = field(init=False, factory=list)
