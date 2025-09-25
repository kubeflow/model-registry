"""Experiment types for model registry.

Experiment types are used to map between proto types based on experiments and Python objects.


"""

from __future__ import annotations  # noqa: I001

from typing import Any
from typing_extensions import override

from mr_openapi import (
    Experiment as ExperimentBaseModel,
)
from mr_openapi import (
    ExperimentCreate,
    ExperimentRunCreate,
    ExperimentRunState,
    ExperimentRunUpdate,
    ExperimentState,
    ExperimentUpdate,
)
from mr_openapi import (
    ExperimentRun as ExperimentRunBaseModel,
)

from .base import BaseResourceModel


class Experiment(BaseResourceModel):
    """Represents an experiment model.

    Attributes:
        name: Name of the experiment.
        owner: Owner of the experiment.
        description: Description of the experiment.
        external_id: External ID of the experiment.
        state: State of the experiment.
        custom_properties: Custom properties (metadata)of the experiment.
    """

    name: str
    owner: str | None = None
    description: str | None = None
    external_id: str | None = None
    state: ExperimentState | None = None
    custom_properties: dict[str, Any] | None = None

    @override
    def create(self, **kwargs) -> ExperimentCreate:
        return ExperimentCreate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "custom_properties")),
            **kwargs,
        )

    @override
    def update(self, **kwargs) -> ExperimentUpdate:
        return ExperimentUpdate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "name", "custom_properties")),
            **kwargs,
        )

    @classmethod
    @override
    def from_basemodel(cls, source: ExperimentBaseModel) -> Experiment:
        assert source.name
        return cls(
            id=source.id,
            name=source.name,
            owner=source.owner,
            description=source.description,
            external_id=source.external_id,
            create_time_since_epoch=source.create_time_since_epoch,
            last_update_time_since_epoch=source.last_update_time_since_epoch,
            state=source.state,
            custom_properties=cls._unmap_custom_properties(source.custom_properties)
            if source.custom_properties
            else None,
        )


class ExperimentRun(BaseResourceModel):
    """Represents an experiment run model.

    Attributes:
        name: Name of the experiment run.
        owner: Owner of the experiment run.
        description: Description of the experiment run.
        external_id: External ID of the experiment run.
        state: State of the experiment.
        custom_properties: Custom properties (metadata)of the experiment.
    """

    name: str
    experiment_id: str
    owner: str | None = None
    description: str | None = None
    external_id: str | None = None
    state: ExperimentRunState | None = None
    custom_properties: dict[str, Any] | None = None

    @override
    def create(self, **kwargs) -> ExperimentRunCreate:
        return ExperimentRunCreate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "custom_properties")),
            **kwargs,
        )

    @override
    def update(self, **kwargs) -> ExperimentRunUpdate:
        return ExperimentRunUpdate(
            customProperties=self._map_custom_properties(),
            **self._props_as_dict(exclude=("id", "name", "custom_properties")),
            **kwargs,
        )

    @classmethod
    @override
    def from_basemodel(cls, source: ExperimentRunBaseModel) -> ExperimentRun:
        assert source.name
        return cls(
            id=source.id,
            name=source.name,
            experiment_id=source.experiment_id,
            owner=source.owner,
            description=source.description,
            external_id=source.external_id,
            create_time_since_epoch=source.create_time_since_epoch,
            last_update_time_since_epoch=source.last_update_time_since_epoch,
            state=source.state,
            custom_properties=cls._unmap_custom_properties(source.custom_properties)
            if source.custom_properties
            else None,
        )
