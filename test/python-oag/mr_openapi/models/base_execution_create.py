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


from typing import Dict, Optional
from pydantic import BaseModel, Field, StrictStr
from mr_openapi.models.execution_state import ExecutionState
from mr_openapi.models.metadata_value import MetadataValue

class BaseExecutionCreate(BaseModel):
    """
    BaseExecutionCreate
    """
    last_known_state: Optional[ExecutionState] = Field(None, alias="lastKnownState")
    custom_properties: Optional[Dict[str, MetadataValue]] = Field(None, alias="customProperties", description="User provided custom properties which are not defined by its type.")
    description: Optional[StrictStr] = Field(None, description="An optional description about the resource.")
    external_id: Optional[StrictStr] = Field(None, alias="externalId", description="The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance.")
    name: Optional[StrictStr] = Field(None, description="The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set.")
    __properties = ["lastKnownState", "customProperties", "description", "externalId", "name"]

    class Config:
        """Pydantic configuration"""
        allow_population_by_field_name = True
        validate_assignment = True

    def to_str(self) -> str:
        """Returns the string representation of the model using alias"""
        return pprint.pformat(self.dict(by_alias=True))

    def to_json(self) -> str:
        """Returns the JSON representation of the model using alias"""
        return json.dumps(self.to_dict())

    @classmethod
    def from_json(cls, json_str: str) -> BaseExecutionCreate:
        """Create an instance of BaseExecutionCreate from a JSON string"""
        return cls.from_dict(json.loads(json_str))

    def to_dict(self):
        """Returns the dictionary representation of the model using alias"""
        _dict = self.dict(by_alias=True,
                          exclude={
                          },
                          exclude_none=True)
        # override the default output from pydantic by calling `to_dict()` of each value in custom_properties (dict)
        _field_dict = {}
        if self.custom_properties:
            for _key in self.custom_properties:
                if self.custom_properties[_key]:
                    _field_dict[_key] = self.custom_properties[_key].to_dict()
            _dict['customProperties'] = _field_dict
        return _dict

    @classmethod
    def from_dict(cls, obj: dict) -> BaseExecutionCreate:
        """Create an instance of BaseExecutionCreate from a dict"""
        if obj is None:
            return None

        if not isinstance(obj, dict):
            return BaseExecutionCreate.parse_obj(obj)

        _obj = BaseExecutionCreate.parse_obj({
            "last_known_state": obj.get("lastKnownState"),
            "custom_properties": dict(
                (_k, MetadataValue.from_dict(_v))
                for _k, _v in obj.get("customProperties").items()
            )
            if obj.get("customProperties") is not None
            else None,
            "description": obj.get("description"),
            "external_id": obj.get("externalId"),
            "name": obj.get("name")
        })
        return _obj


