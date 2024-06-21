"""Context types for model registry.

Contexts group related Artifacts together.
They provide a way to organize and categorize components in a workflow.

Those types are used to map between proto types based on contexts and Python objects.

Todo:
    * Move part of the description to API Reference docs (#120).
"""

from __future__ import annotations

from .base import BaseResourceModel


class ModelVersion(BaseResourceModel):
    """Represents a model version.

    Attributes:
        name: Name of this version.
        author: Author of the model version.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
        artifacts: Artifacts associated with this version.
    """

    author: str
    # state: ModelVersionState = ModelVersionState.LIVE


class RegisteredModel(BaseResourceModel):
    """Represents a registered model.

    Attributes:
        name: Registered model name.
        owner: Owner of this Registered Model.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
    """

    owner: str | None = None
    # state: RegisteredModelState = RegisteredModelState.LIVE
