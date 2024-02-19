"""Tests for artifact type mapping.

Todo:
    * should we parametrize the tests?
"""

import pytest
from ml_metadata.proto import Artifact
from model_registry.types import ModelArtifact

from .. import Mapped


@pytest.fixture()
def complete_model() -> Mapped:
    proto_model = Artifact()
    proto_model.name = "test_prefix:test_model"
    proto_model.type_id = 1
    proto_model.external_id = "test_external_id"
    proto_model.state = Artifact.UNKNOWN
    proto_model.uri = "test_uri"
    proto_model.properties["description"].string_value = "test description"
    proto_model.properties["model_format_name"].string_value = "test_format"
    proto_model.properties["model_format_version"].string_value = "test_format_version"
    proto_model.properties["storage_key"].string_value = "test_storage_key"
    proto_model.properties["storage_path"].string_value = "test_storage_path"
    proto_model.properties["service_account_name"].string_value = "test_account_name"

    py_model = ModelArtifact(
        "test_model",
        "test_uri",
        description="test description",
        external_id="test_external_id",
        model_format_name="test_format",
        model_format_version="test_format_version",
        storage_key="test_storage_key",
        storage_path="test_storage_path",
        service_account_name="test_account_name",
    )

    return Mapped(proto_model, py_model)


@pytest.fixture()
def minimal_model() -> Mapped:
    proto_model = Artifact()
    proto_model.name = "test_prefix:test_model"
    proto_model.type_id = 1
    proto_model.state = Artifact.UNKNOWN
    proto_model.uri = "test_uri"

    py_model = ModelArtifact("test_model", "test_uri")
    return Mapped(proto_model, py_model)


def test_partial_model_mapping(monkeypatch, minimal_model: Mapped):
    monkeypatch.setattr(ModelArtifact, "mlmd_name_prefix", "test_prefix")

    mapped_model = minimal_model.py.map(1)
    proto_model = minimal_model.proto
    assert mapped_model.name == proto_model.name
    assert mapped_model.type_id == proto_model.type_id
    assert mapped_model.state == proto_model.state
    assert mapped_model.uri == proto_model.uri


def test_full_model_mapping(monkeypatch, complete_model: Mapped):
    monkeypatch.setattr(ModelArtifact, "mlmd_name_prefix", "test_prefix")

    mapped_model = complete_model.py.map(1)
    proto_model = complete_model.proto
    assert mapped_model.name == proto_model.name
    assert mapped_model.type_id == proto_model.type_id
    assert mapped_model.state == proto_model.state
    assert mapped_model.uri == proto_model.uri
    assert mapped_model.external_id == proto_model.external_id
    assert mapped_model.properties == proto_model.properties


def test_partial_model_unmapping(minimal_model: Mapped):
    unmapped_model = ModelArtifact.unmap(minimal_model.proto)
    py_model = minimal_model.py
    assert unmapped_model.name == py_model.name
    assert unmapped_model.state == py_model.state
    assert unmapped_model.uri == py_model.uri


def test_full_model_unmapping(complete_model: Mapped):
    unmapped_model = ModelArtifact.unmap(complete_model.proto)
    py_model = complete_model.py
    assert unmapped_model.name == py_model.name
    assert unmapped_model.state == py_model.state
    assert unmapped_model.uri == py_model.uri
    assert unmapped_model.external_id == py_model.external_id
    assert unmapped_model.description == py_model.description
    assert unmapped_model.model_format_name == py_model.model_format_name
    assert unmapped_model.model_format_version == py_model.model_format_version
    assert unmapped_model.storage_key == py_model.storage_key
    assert unmapped_model.storage_path == py_model.storage_path
    assert unmapped_model.service_account_name == py_model.service_account_name


def test_model_prefix_generation(minimal_model: Mapped):
    name1 = minimal_model.py.map(1).name
    name2 = minimal_model.py.map(1).name
    assert name1 != name2
