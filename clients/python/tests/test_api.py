import secrets
from typing import Callable

import pytest
import requests
import schemathesis
from hypothesis import HealthCheck, settings

from .conftest import REGISTRY_URL
from .constants import ARTIFACT_STATES, ARTIFACT_TYPE_PARAMS, DEFAULT_API_TIMEOUT

schema = schemathesis.pytest.from_fixture("generated_schema")


schema = (
    schema
    .exclude(
        path="/api/model_registry/v1alpha3/artifacts/{id}",
        method="PATCH"
    )
    .exclude(
        path="/api/model_registry/v1alpha3/model_versions/{modelversionId}/artifacts",
        method="POST"
    )
)
@schema.parametrize()
@settings(
    max_examples=100,
    deadline=None,
    suppress_health_check=[
        HealthCheck.filter_too_much,
        HealthCheck.too_slow,
        HealthCheck.data_too_large,
    ],
)
@pytest.mark.fuzz
def test_mr_api_stateless(auth_headers: dict, case: schemathesis.Case):
    """Test the Model Registry API endpoints.

    This test uses schemathesis to generate and validate API requests
    """

    case.call_and_validate(headers=auth_headers)

@pytest.mark.fuzz
@pytest.mark.parametrize(("artifact_type", "uri_prefix"), ARTIFACT_TYPE_PARAMS)
@pytest.mark.parametrize("state", ARTIFACT_STATES)
def test_post_model_version_artifacts(auth_headers: dict, artifact_type: str, uri_prefix: str, state: str, cleanup_artifacts: Callable):
    """
    Direct test for POST /api/model_registry/v1alpha3/model_versions/{modelversionId}/artifacts.
    """
    model_version_id = str(secrets.randbelow(2000000000 - 100000 + 1) + 100000)

    endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/model_versions/{model_version_id}/artifacts"

    payload = {
        "artifactType": artifact_type,
        "name": "my-test-model-artifact-post",
        "uri": f"{uri_prefix}my-test-model.pkl",
        "state": state,
        "description": "A test model artifact created via direct POST test.",
        "externalId": str(secrets.randbelow(2000000000 - 100000 + 1) + 100000)
    }

    response = requests.post(endpoint, headers=auth_headers, json=payload, timeout=DEFAULT_API_TIMEOUT)
    artifact_id = response.json()["id"]
    cleanup_artifacts(artifact_id)

    assert response.status_code in {200, 201}, f"Expected 200 or 201, got {response.status_code}: {response.text}"
    response_json = response.json()
    assert response_json.get("id"), "Response body should contain 'id'"
    assert response_json.get("name") == payload["name"], "Response name should match payload name"
    assert response_json.get("artifactType") == payload["artifactType"], "Response artifactType should match payload"


@pytest.mark.fuzz
@pytest.mark.parametrize(("artifact_type", "uri_prefix"), ARTIFACT_TYPE_PARAMS)
def test_patch_artifact(auth_headers: dict, artifact_resource: Callable, artifact_type: str, uri_prefix: str):
    """
    Direct test for PATCH /api/model_registry/v1alpha3/artifacts/{id}.
    """
    initial_state = "PENDING"
    target_state = "LIVE"

    create_payload = {
        "artifactType": artifact_type,
        "name": "test-create-for-patch",
        "uri": "s3://my-test-bucket/models/initial-model.pkl",
        "state": initial_state,
    }
    if artifact_type == "model-artifact":
        create_payload["modelFormatName"] = "tensorflow"
        create_payload["modelFormatVersion"] = "1.0"


    with artifact_resource(auth_headers, create_payload) as artifact_id:
        patch_endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
        patch_payload = {
            "artifactType": artifact_type,
            "description": f"Updated description for {artifact_type} ({target_state})",
            "state": target_state,
        }
        patch_response = requests.patch(patch_endpoint, headers=auth_headers, json=patch_payload, timeout=DEFAULT_API_TIMEOUT)
        assert patch_response.status_code == 200
        patch_response_json = patch_response.json()
        assert patch_response_json.get("id") == artifact_id
        assert patch_response_json.get("description") == patch_payload["description"]
        assert patch_response_json.get("state") == patch_payload["state"]
        assert patch_response_json.get("artifactType") == artifact_type
