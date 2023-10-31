"""Base types for model registry."""

from __future__ import annotations

from abc import ABC, abstractmethod
from collections.abc import Mapping
from inspect import isabstract
from typing import Any, ClassVar, Optional

from attrs import define, field

from model_registry.store import ProtoType, ScalarType


class Mappable(ABC):
    """Interface for types that can be mapped to and from proto types."""

    @classmethod
    def get_proto_type_name(cls) -> str:
        """Name of the proto type

        Returns:
            str: Name of the class prefixed with `odh.`
        """
        return f"odh.{cls.__name__}"

    @property
    @abstractmethod
    def proto_name(self) -> str:
        """Name of the proto object.
        """
        pass

    @abstractmethod
    def map(self) -> ProtoType:
        """Map to a proto object.

        Returns:
            ProtoType: Proto object.
        """
        pass

    @classmethod
    @abstractmethod
    def unmap(cls, mlmd_obj: ProtoType) -> Mappable:
        """Map from a proto object.

        Args:
            ProtoType: Proto object.

        Returns:
            Mappable: Python object.
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
        """Prefix to be used in the proto object.
        """
        pass


@define(slots=False, init=False)
class ProtoBase(Mappable, ABC):
    """Abstract base type for protos

    This is a type defining common functionality for all types representing Model Registry protos,
    such as Artifacts, Contexts, and Executions.

    Attributes:
        id (str): Protobuf object ID. Auto-assigned when put on the server.
        name (str): Name of the object.
        description (str, optional): Description of the object.
        external_id (str, optional): Customizable ID. Has to be unique among instances of the same type.
        create_time_since_epoch (int): Seconds elapsed since object creation time, measured against epoch.
        last_update_time_since_epoch (int): Seconds elapsed since object last update time, measured against epoch.
    """

    name: str = field(init=False)
    id: Optional[str] = field(init=False, default=None)
    description: Optional[str] = field(kw_only=True, default=None)
    external_id: Optional[str] = field(kw_only=True, default=None)
    create_time_since_epoch: Optional[int] = field(init=False, default=None)
    last_update_time_since_epoch: Optional[int] = field(init=False, default=None)

    _types_map: ClassVar[dict[str, ProtoBase]] = {}

    def __init_subclass__(cls):
        super().__init_subclass__()
        if isabstract(cls):
            return
        cls._types_map[cls.get_proto_type_name()] = cls

    @classmethod
    def get_subclass(cls, proto_type_name: str) -> Mappable:
        """Get a subclass by proto type name.

        Args:
            proto_type_name (str): Name of the proto type.

        Returns:
            Mappable: Subclass.
        """
        return cls._types_map[proto_type_name]

    @property
    def proto_name(self) -> str:
        if isinstance(self, Prefixable):
            return f"{self.mlmd_name_prefix}:{self.name}"
        return self.name

    @classmethod
    @abstractmethod
    def get_proto_type(cls) -> type[ProtoType]:
        """Proto type associated with this class.

        Returns:
            ProtoType: Proto type.
        """
        pass

    @staticmethod
    def _map_props(
        py_props: Mapping[str, Optional[ScalarType]], mlmd_props: dict[str, Any]
    ):
        """Map properties from Python to proto.

        Args:
            py_props (dict[str, ScalarType]): Python properties.
            mlmd_props (dict[str, Any]): Proto properties, will be modified in place.
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
                raise Exception(f"Unsupported type: {type(value)}")

    def map(self) -> ProtoType:
        mlmd_obj = (self.get_proto_type())()
        mlmd_obj.name = self.proto_name
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
            mlmd_props (dict[str, Any]): Proto properties.

        Returns:
            dict[str, ScalarType]: Python properties.
        """
        py_props: dict[str, ScalarType] = {}
        for key, value in mlmd_props.items():
            # TODO: use pattern matching here (3.10)
            if value.HasField("bool_value"):
                py_props[key] = value.bool_value
            elif value.HasField("int_value"):
                py_props[key] = value.int_value
            elif value.HasField("double_value"):
                py_props[key] = value.double_value
            elif value.HasField("string_value"):
                py_props[key] = value.string_value
            else:
                raise Exception(f"Unsupported type: {type(value)}")

        return py_props

    @classmethod
    def unmap(cls, mlmd_obj: ProtoType) -> ProtoBase:
        py_obj = cls.__new__(cls)
        py_obj.id = str(mlmd_obj.id)
        if isinstance(py_obj, Prefixable):
            name: str = mlmd_obj.name
            assert ':' in name, f"Expected {name} to be prefixed"
            py_obj.name = name.split(':', 1)[1]
        else:
            py_obj.name = mlmd_obj.name
        py_obj.description = mlmd_obj.properties["description"].string_value
        py_obj.external_id = mlmd_obj.external_id
        py_obj.create_time_since_epoch = mlmd_obj.create_time_since_epoch
        py_obj.last_update_time_since_epoch = mlmd_obj.last_update_time_since_epoch
        return py_obj
