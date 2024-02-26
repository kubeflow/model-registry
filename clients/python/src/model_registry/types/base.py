"""Base types for model registry."""

from __future__ import annotations

from abc import ABC, abstractmethod
from collections.abc import Mapping
from typing import Any, TypeVar, get_args

from attrs import define, field
from typing_extensions import override

from model_registry.store import ProtoType, ScalarType


class Mappable(ABC):
    """Interface for types that can be mapped to and from proto types."""

    @classmethod
    def get_proto_type_name(cls) -> str:
        """Name of the proto type.

        Returns:
            Name of the class prefixed with `kf.`
        """
        return f"kf.{cls.__name__}"

    @property
    @abstractmethod
    def proto_name(self) -> str:
        """Name of the proto object."""
        pass

    @abstractmethod
    def map(self, type_id: int) -> ProtoType:
        """Map to a proto object.

        Args:
            type_id (int): ID of the type.

        Returns:
            Proto object.
        """
        pass

    @classmethod
    @abstractmethod
    def unmap(cls, mlmd_obj: ProtoType) -> Mappable:
        """Map from a proto object.

        Args:
            mlmd_obj: Proto object.

        Returns:
            Python object.
        """
        pass


class Prefixable(ABC):
    """Interface for types that are prefixed.

    We use prefixes to ensure that the user can insert more than one instance of the same type
    with the same name/external_id.
    """

    @property
    @abstractmethod
    def mlmd_name_prefix(self) -> str:
        """Prefix to be used in the proto object."""
        pass


@define(slots=False, init=False)
class ProtoBase(Mappable, ABC):
    """Abstract base type for protos.

    This is a type defining common functionality for all types representing Model Registry protos,
    such as Artifacts, Contexts, and Executions.

    Attributes:
        id: Protobuf object ID. Auto-assigned when put on the server.
        name: Name of the object.
        description: Description of the object.
        external_id: Customizable ID. Has to be unique among instances of the same type.
        create_time_since_epoch: Seconds elapsed since object creation time, measured against epoch.
        last_update_time_since_epoch: Seconds elapsed since object last update time, measured against epoch.
    """

    name: str = field(init=False)
    id: str | None = field(init=False, default=None)
    description: str | None = field(kw_only=True, default=None)
    external_id: str | None = field(kw_only=True, default=None)
    create_time_since_epoch: int | None = field(init=False, default=None)
    last_update_time_since_epoch: int | None = field(init=False, default=None)

    @property
    @override
    def proto_name(self) -> str:
        if isinstance(self, Prefixable):
            return f"{self.mlmd_name_prefix}:{self.name}"
        return self.name

    @classmethod
    @abstractmethod
    def get_proto_type(cls) -> type[ProtoType]:
        """Proto type associated with this class.

        Returns:
            Proto type.
        """
        pass

    @staticmethod
    def _map_props(
        py_props: Mapping[str, ScalarType | None], mlmd_props: dict[str, Any]
    ):
        """Map properties from Python to proto.

        Args:
            py_props: Python properties.
            mlmd_props: Proto properties, will be modified in place.
        """
        for key, value in py_props.items():
            if value is None:
                continue
            # TODO: use pattern matching here (3.10)
            if isinstance(value, bool):
                mlmd_props[key].bool_value = value
            elif isinstance(value, int):
                mlmd_props[key].int_value = value
            elif isinstance(value, float):
                mlmd_props[key].double_value = value
            elif isinstance(value, str):
                mlmd_props[key].string_value = value
            else:
                msg = f"Unsupported type: {type(value)}"
                raise Exception(msg)

    @override
    def map(self, type_id: int) -> ProtoType:
        mlmd_obj = (self.get_proto_type())()
        mlmd_obj.name = self.proto_name
        mlmd_obj.type_id = type_id
        if self.id:
            mlmd_obj.id = int(self.id)
        if self.external_id:
            mlmd_obj.external_id = self.external_id
        if self.description:
            mlmd_obj.properties["description"].string_value = self.description
        return mlmd_obj

    @staticmethod
    def _unmap_props(mlmd_props: dict[str, Any]) -> dict[str, ScalarType]:
        """Map properties from proto to Python.

        Args:
            mlmd_props: Proto properties.

        Returns:
            Python properties.
        """
        py_props: dict[str, ScalarType] = {}
        for key, prop in mlmd_props.items():
            _, value = prop.ListFields()[0]
            if not isinstance(value, get_args(ScalarType)):
                msg = f"Unsupported type {type(value)} on key {key}"
                raise Exception(msg)
            py_props[key] = value

        return py_props

    T = TypeVar("T", bound="ProtoBase")

    @classmethod
    @override
    def unmap(cls: type[T], mlmd_obj: ProtoType) -> T:
        py_obj = cls.__new__(cls)
        py_obj.id = str(mlmd_obj.id)
        if isinstance(py_obj, Prefixable):
            name: str = mlmd_obj.name
            assert ":" in name, f"Expected {name} to be prefixed"
            py_obj.name = name.split(":", 1)[1]
        else:
            py_obj.name = mlmd_obj.name
        py_obj.description = mlmd_obj.properties["description"].string_value
        py_obj.external_id = mlmd_obj.external_id
        py_obj.create_time_since_epoch = mlmd_obj.create_time_since_epoch
        py_obj.last_update_time_since_epoch = mlmd_obj.last_update_time_since_epoch
        return py_obj
