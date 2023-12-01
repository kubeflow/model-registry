"""Tests for context type mapping.

Todo:
    * should we parametrize the tests?
"""

import pytest
from ml_metadata.proto import Context
from model_registry.types import ContextState, ModelVersion

from .. import Mapped


@pytest.fixture()
def full_model_version() -> Mapped:
    proto_version = Context()
    proto_version.name = "1:1.0.0"
    proto_version.type_id = 2
    proto_version.external_id = "test_external_id"
    proto_version.properties["description"].string_value = "test description"
    proto_version.properties["model_name"].string_value = "test_model"
    proto_version.properties["author"].string_value = "test_author"
    proto_version.properties["tags"].string_value = "test_tag1,test_tag2"
    proto_version.properties["state"].string_value = "ARCHIVED"
    proto_version.custom_properties["int_key"].int_value = 1
    proto_version.custom_properties["float_key"].double_value = 1.0
    proto_version.custom_properties["bool_key"].bool_value = True
    proto_version.custom_properties["str_key"].string_value = "test_str"

    py_version = ModelVersion(
        "test_model",
        "1.0.0",
        "test_author",
        external_id="test_external_id",
        description="test description",
    )
    py_version._registered_model_id = 1
    py_version.tags = ["test_tag1", "test_tag2"]
    py_version.metadata = {
        "int_key": 1,
        "float_key": 1.0,
        "bool_key": True,
        "str_key": "test_str",
    }
    py_version.state = ContextState.ARCHIVED
    return Mapped(proto_version, py_version)


@pytest.fixture()
def minimal_model_version() -> Mapped:
    proto_version = Context()
    proto_version.name = "1:1.0.0"
    proto_version.type_id = 2
    proto_version.properties["model_name"].string_value = "test_model"
    proto_version.properties["author"].string_value = "test_author"
    proto_version.properties["state"].string_value = "LIVE"

    py_version = ModelVersion("test_model", "1.0.0", "test_author")
    py_version._registered_model_id = 1
    return Mapped(proto_version, py_version)


def test_partial_model_version_mapping(minimal_model_version: Mapped):
    mapped_version = minimal_model_version.py.map(2)
    proto_version = minimal_model_version.proto
    assert mapped_version.name == proto_version.name
    assert mapped_version.type_id == proto_version.type_id
    assert mapped_version.properties == proto_version.properties


def test_partial_model_version_unmapping(minimal_model_version: Mapped):
    unmapped_version = ModelVersion.unmap(minimal_model_version.proto)
    py_version = minimal_model_version.py
    assert unmapped_version.version == py_version.version
    assert unmapped_version.model_name == py_version.model_name
    assert unmapped_version.author == py_version.author
    assert unmapped_version.state == py_version.state
    assert unmapped_version.tags == py_version.tags
    assert unmapped_version.metadata == py_version.metadata


def test_full_model_version_mapping(full_model_version: Mapped):
    mapped_version = full_model_version.py.map(2)
    proto_version = full_model_version.proto
    assert mapped_version.name == proto_version.name
    assert mapped_version.type_id == proto_version.type_id
    assert mapped_version.external_id == proto_version.external_id
    assert mapped_version.properties == proto_version.properties
    assert mapped_version.custom_properties == proto_version.custom_properties


def test_full_model_version_unmapping(full_model_version: Mapped):
    unmapped_version = ModelVersion.unmap(full_model_version.proto)
    py_version = full_model_version.py
    assert unmapped_version.version == py_version.version
    assert unmapped_version.description == py_version.description
    assert unmapped_version.external_id == py_version.external_id
    assert unmapped_version.model_name == py_version.model_name
    assert unmapped_version.author == py_version.author
    assert unmapped_version.state == py_version.state
    assert unmapped_version.tags == py_version.tags
    assert unmapped_version.metadata == py_version.metadata
