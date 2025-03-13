# coding: utf-8

"""
    Model Registry REST API

    REST API for Model Registry to create and manage ML model metadata

    The version of the OpenAPI document: v1alpha3
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


from __future__ import annotations
import pprint
import re  # noqa: F401
import json

from pydantic import BaseModel, ConfigDict, Field, StrictStr
from typing import Any, ClassVar, Dict, List, Optional
from mr_openapi.models.metadata_value import MetadataValue
from mr_openapi.models.model_version_state import ModelVersionState
from typing import Optional, Set
from typing_extensions import Self

class ModelVersion(BaseModel):
    """
    Represents a ModelVersion belonging to a RegisteredModel.
    """ # noqa: E501
    custom_properties: Optional[Dict[str, MetadataValue]] = Field(default=None, description="User provided custom properties which are not defined by its type.", alias="customProperties")
    description: Optional[StrictStr] = Field(default=None, description="An optional description about the resource.")
    external_id: Optional[StrictStr] = Field(default=None, description="The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance.", alias="externalId")
    name: StrictStr = Field(description="The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set.")
    state: Optional[ModelVersionState] = None
    author: Optional[StrictStr] = Field(default=None, description="Name of the author.")
    registered_model_id: StrictStr = Field(description="ID of the `RegisteredModel` to which this version belongs.", alias="registeredModelId")
    id: Optional[StrictStr] = Field(default=None, description="The unique server generated id of the resource.")
    create_time_since_epoch: Optional[StrictStr] = Field(default=None, description="Output only. Create time of the resource in millisecond since epoch.", alias="createTimeSinceEpoch")
    last_update_time_since_epoch: Optional[StrictStr] = Field(default=None, description="Output only. Last update time of the resource since epoch in millisecond since epoch.", alias="lastUpdateTimeSinceEpoch")
    __properties: ClassVar[List[str]] = ["customProperties", "description", "externalId", "name", "state", "author", "registeredModelId", "id", "createTimeSinceEpoch", "lastUpdateTimeSinceEpoch"]

    model_config = ConfigDict(
        populate_by_name=True,
        validate_assignment=True,
        protected_namespaces=(),
    )


    def to_str(self) -> str:
        """Returns the string representation of the model using alias"""
        return pprint.pformat(self.model_dump(by_alias=True))

    def to_json(self) -> str:
        """Returns the JSON representation of the model using alias"""
        # TODO: pydantic v2: use .model_dump_json(by_alias=True, exclude_unset=True) instead
        return json.dumps(self.to_dict())

    @classmethod
    def from_json(cls, json_str: str) -> Optional[Self]:
        """Create an instance of ModelVersion from a JSON string"""
        return cls.from_dict(json.loads(json_str))

    def to_dict(self) -> Dict[str, Any]:
        """Return the dictionary representation of the model using alias.

        This has the following differences from calling pydantic's
        `self.model_dump(by_alias=True)`:

        * `None` is only added to the output dict for nullable fields that
          were set at model initialization. Other fields with value `None`
          are ignored.
        * OpenAPI `readOnly` fields are excluded.
        * OpenAPI `readOnly` fields are excluded.
        """
        excluded_fields: Set[str] = set([
            "create_time_since_epoch",
            "last_update_time_since_epoch",
        ])

        _dict = self.model_dump(
            by_alias=True,
            exclude=excluded_fields,
            exclude_none=True,
        )
        # override the default output from pydantic by calling `to_dict()` of each value in custom_properties (dict)
        _field_dict = {}
        if self.custom_properties:
            for _key in self.custom_properties:
                if self.custom_properties[_key]:
                    _field_dict[_key] = self.custom_properties[_key].to_dict()
            _dict['customProperties'] = _field_dict
        return _dict

    @classmethod
    def from_dict(cls, obj: Optional[Dict[str, Any]]) -> Optional[Self]:
        """Create an instance of ModelVersion from a dict"""
        if obj is None:
            return None

        if not isinstance(obj, dict):
            return cls.model_validate(obj)

        _obj = cls.model_validate({
            "customProperties": dict(
                (_k, MetadataValue.from_dict(_v))
                for _k, _v in obj["customProperties"].items()
            )
            if obj.get("customProperties") is not None
            else None,
            "description": obj.get("description"),
            "externalId": obj.get("externalId"),
            "name": obj.get("name"),
            "state": obj.get("state"),
            "author": obj.get("author"),
            "registeredModelId": obj.get("registeredModelId"),
            "id": obj.get("id"),
            "createTimeSinceEpoch": obj.get("createTimeSinceEpoch"),
            "lastUpdateTimeSinceEpoch": obj.get("lastUpdateTimeSinceEpoch")
        })
        return _obj


