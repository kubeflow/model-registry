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


from typing import List, Optional
from pydantic import BaseModel, Field, StrictInt, StrictStr, conlist
from mr_openapi.models.model_version import ModelVersion

class ModelVersionList(BaseModel):
    """
    List of ModelVersion entities.  # noqa: E501
    """
    next_page_token: StrictStr = Field(..., alias="nextPageToken", description="Token to use to retrieve next page of results.")
    page_size: StrictInt = Field(..., alias="pageSize", description="Maximum number of resources to return in the result.")
    size: StrictInt = Field(..., description="Number of items in result list.")
    items: Optional[conlist(ModelVersion)] = Field(None, description="Array of `ModelVersion` entities.")
    __properties = ["nextPageToken", "pageSize", "size", "items"]

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
    def from_json(cls, json_str: str) -> ModelVersionList:
        """Create an instance of ModelVersionList from a JSON string"""
        return cls.from_dict(json.loads(json_str))

    def to_dict(self):
        """Returns the dictionary representation of the model using alias"""
        _dict = self.dict(by_alias=True,
                          exclude={
                          },
                          exclude_none=True)
        # override the default output from pydantic by calling `to_dict()` of each item in items (list)
        _items = []
        if self.items:
            for _item in self.items:
                if _item:
                    _items.append(_item.to_dict())
            _dict['items'] = _items
        return _dict

    @classmethod
    def from_dict(cls, obj: dict) -> ModelVersionList:
        """Create an instance of ModelVersionList from a dict"""
        if obj is None:
            return None

        if not isinstance(obj, dict):
            return ModelVersionList.parse_obj(obj)

        _obj = ModelVersionList.parse_obj({
            "next_page_token": obj.get("nextPageToken"),
            "page_size": obj.get("pageSize"),
            "size": obj.get("size"),
            "items": [ModelVersion.from_dict(_item) for _item in obj.get("items")] if obj.get("items") is not None else None
        })
        return _obj


