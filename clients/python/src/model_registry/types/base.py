"""Base type for protos"""

from abc import ABC
from typing import Optional

from attrs import define, field


@define(slots=False, init=False)
class ProtoBase(ABC):
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
    description: str = field(kw_only=True)
    external_id: Optional[str] = field(kw_only=True, default=None)
    create_time_since_epoch: Optional[int] = field(init=False, default=None)
    last_update_time_since_epoch: Optional[int] = field(init=False, default=None)
