import logging
import secrets
import time
from typing import Any, Callable

import pytest
import requests
import schemathesis
from hypothesis import HealthCheck, settings

from .conftest import REGISTRY_URL
from .constants import ARTIFACT_STATES, ARTIFACT_TYPE_PARAMS, DEFAULT_API_TIMEOUT


# Helper functions for common operations
def generate_random_id() -> str:
    """Generate a random ID between 100000 and 2000000000."""
    return str(secrets.randbelow(2000000000 - 100000 + 1) + 100000)


def generate_unique_timestamp() -> str:
    """Generate a unique timestamp in milliseconds with random offset."""
    current_ms = int(time.time() * 1000)
    random_offset = secrets.randbelow(1000000)  # Random offset up to 1 million ms (1000 seconds)
    return str(current_ms + random_offset)


def build_artifact_payload(artifact_type: str, uri_prefix: str, state: str, name: str,
                          description: str = None, external_id: str = None) -> dict[str, Any]:
    """Build a payload for creating an artifact based on its type.

    Args:
        artifact_type: Type of artifact (e.g., 'model-artifact', 'doc-artifact', etc.)
        uri_prefix: Prefix for URI (e.g., 's3://', 'https://')
        state: State of the artifact
        name: Name of the artifact
        description: Optional description
        external_id: Optional external ID, generates random if not provided

    Returns:
        Dictionary containing the artifact payload
    """
    payload = {
        "artifactType": artifact_type,
        "name": name,
        "state": state,
        "externalId": external_id or generate_random_id()
    }

    if description:
        payload["description"] = description

    # Add type-specific properties
    if artifact_type == "model-artifact":
        payload["modelFormatName"] = "tensorflow"
        payload["modelFormatVersion"] = "1.0"
        payload["uri"] = f"{uri_prefix}my-test-model.pkl"
    elif artifact_type == "doc-artifact":
        payload["uri"] = f"{uri_prefix}documentation.pdf"
    elif artifact_type == "dataset-artifact":
        payload["uri"] = f"{uri_prefix}dataset.parquet"
        payload["sourceType"] = "s3"
        payload["source"] = "s3://test-bucket/datasets/"
    elif artifact_type == "metric":
        payload["value"] = 0.95  # Example accuracy value
        payload["timestamp"] = generate_unique_timestamp()
    elif artifact_type == "parameter":
        payload["value"] = "0.001"  # Example learning rate
        payload["parameterType"] = "string"

    return payload


def validate_artifact_response(response: requests.Response, expected_payload: dict[str, Any]) -> str:
    """Validate artifact creation response and return the artifact ID.

    Args:
        response: The HTTP response from the API
        expected_payload: The payload that was sent, used for validation

    Returns:
        The artifact ID from the response

    Raises:
        AssertionError: If validation fails
    """
    # Check response status
    assert response.status_code in {200, 201}, f"Expected 200 or 201, got {response.status_code}: {response.text}"

    response_json = response.json()
    assert response_json.get("id"), "Response body should contain 'id'"

    # Validate response matches payload
    assert response_json.get("name") == expected_payload["name"], "Response name should match payload name"
    assert response_json.get("artifactType") == expected_payload["artifactType"], "Response artifactType should match payload"

    return response_json["id"]


def create_experiment_and_run(auth_headers: dict[str, str]) -> tuple[str, str]:
    """Create an experiment and an experiment run.

    Args:
        auth_headers: Authentication headers for the API

    Returns:
        Tuple of (experiment_id, experiment_run_id)

    Raises:
        AssertionError: If creation fails
    """
    # Create an experiment
    experiment_payload = {
        "name": f"test-experiment-{secrets.randbelow(1000000)}",
        "externalId": generate_random_id(),
        "description": "Test experiment for artifact testing"
    }
    exp_response = requests.post(
        f"{REGISTRY_URL}/api/model_registry/v1alpha3/experiments",
        headers=auth_headers,
        json=experiment_payload,
        timeout=DEFAULT_API_TIMEOUT
    )
    assert exp_response.status_code in {200, 201}, f"Failed to create experiment: {exp_response.text}"
    experiment_id = exp_response.json()["id"]

    # Create an experiment run
    experiment_run_payload = {
        "experimentId": experiment_id,
        "name": f"test-run-{secrets.randbelow(1000000)}",
        "externalId": generate_random_id(),
        "description": "Test experiment run for artifact testing"
    }
    run_response = requests.post(
        f"{REGISTRY_URL}/api/model_registry/v1alpha3/experiment_runs",
        headers=auth_headers,
        json=experiment_run_payload,
        timeout=DEFAULT_API_TIMEOUT
    )
    assert run_response.status_code in {200, 201}, f"Failed to create experiment run: {run_response.text}"
    experiment_run_id = run_response.json()["id"]

    return experiment_id, experiment_run_id


def cleanup_experiment_and_run(auth_headers: dict[str, str], experiment_id: str, experiment_run_id: str) -> None:
    """Best effort cleanup of experiment run and experiment.

    Args:
        auth_headers: Authentication headers for the API
        experiment_id: ID of the experiment to delete
        experiment_run_id: ID of the experiment run to delete
    """
    try:
        requests.delete(
            f"{REGISTRY_URL}/api/model_registry/v1alpha3/experiment_runs/{experiment_run_id}",
            headers=auth_headers,
            timeout=DEFAULT_API_TIMEOUT
        )
        requests.delete(
            f"{REGISTRY_URL}/api/model_registry/v1alpha3/experiments/{experiment_id}",
            headers=auth_headers,
            timeout=DEFAULT_API_TIMEOUT
        )
    except Exception as e:
        logging.warning(f"Failed to cleanup experiment (id={experiment_id}) and/or experiment run (id={experiment_run_id}): {e}")


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
    .exclude(
        path="/api/model_registry/v1alpha3/experiment_runs/{experimentrunId}/artifacts",
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
    model_version_id = generate_random_id()
    endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/model_versions/{model_version_id}/artifacts"

    # Build payload using helper function
    payload = build_artifact_payload(
        artifact_type=artifact_type,
        uri_prefix=uri_prefix,
        state=state,
        name="my-test-model-artifact-post",
        description="A test model artifact created via direct POST test."
    )

    # Make the API request
    response = requests.post(endpoint, headers=auth_headers, json=payload, timeout=DEFAULT_API_TIMEOUT)

    # Validate response and get artifact ID
    artifact_id = validate_artifact_response(response, payload)

    # Cleanup after successful creation
    cleanup_artifacts(artifact_id)


@pytest.mark.fuzz
@pytest.mark.parametrize(("artifact_type", "uri_prefix"), ARTIFACT_TYPE_PARAMS)
@pytest.mark.parametrize("state", ARTIFACT_STATES)
def test_post_experiment_run_artifacts(auth_headers: dict, artifact_type: str, uri_prefix: str, state: str, cleanup_artifacts: Callable):
    """
    Direct test for POST /api/model_registry/v1alpha3/experiment_runs/{experimentrunId}/artifacts.
    """
    # Create experiment and experiment run using helper
    experiment_id, experiment_run_id = create_experiment_and_run(auth_headers)

    endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/experiment_runs/{experiment_run_id}/artifacts"

    # Build payload using helper function
    payload = build_artifact_payload(
        artifact_type=artifact_type,
        uri_prefix=uri_prefix,
        state=state,
        name=f"my-test-experiment-artifact-post-{secrets.randbelow(1000000)}",
        description="A test experiment artifact created via direct POST test."
    )

    # Make the API request
    response = requests.post(endpoint, headers=auth_headers, json=payload, timeout=DEFAULT_API_TIMEOUT)

    # Validate response and get artifact ID
    artifact_id = validate_artifact_response(response, payload)

    # Cleanup artifacts
    cleanup_artifacts(artifact_id)

    # Cleanup experiment and run
    cleanup_experiment_and_run(auth_headers, experiment_id, experiment_run_id)


@pytest.mark.fuzz
@pytest.mark.parametrize(("artifact_type", "uri_prefix"), ARTIFACT_TYPE_PARAMS)
def test_patch_artifact(auth_headers: dict, artifact_resource: Callable, artifact_type: str, uri_prefix: str):
    """
    Direct test for PATCH /api/model_registry/v1alpha3/artifacts/{id}.
    """
    initial_state = "PENDING"
    target_state = "LIVE"

    # Build initial payload for creation
    # Note: Using a hardcoded URI prefix for creation since these are specific test URIs
    build_artifact_payload(
        artifact_type=artifact_type,
        uri_prefix="s3://my-test-bucket/" if artifact_type in ["model-artifact", "dataset-artifact"] else "https://docs.example.com/",
        state=initial_state,
        name="test-create-for-patch"
    )

    # Override some specific properties for the patch test
    create_payload = {
        "artifactType": artifact_type,
        "name": "test-create-for-patch",
        "state": initial_state,
    }

    # Add type-specific properties for creation with test-specific values
    if artifact_type == "model-artifact":
        create_payload["modelFormatName"] = "tensorflow"
        create_payload["modelFormatVersion"] = "1.0"
        create_payload["uri"] = "s3://my-test-bucket/models/initial-model.pkl"
    elif artifact_type == "doc-artifact":
        create_payload["uri"] = "https://docs.example.com/initial-doc.pdf"
    elif artifact_type == "dataset-artifact":
        create_payload["uri"] = "s3://my-test-bucket/datasets/initial-dataset.parquet"
        create_payload["sourceType"] = "s3"
    elif artifact_type == "metric":
        create_payload["value"] = 0.85
        create_payload["timestamp"] = "1000000000"
    elif artifact_type == "parameter":
        create_payload["value"] = "0.01"
        create_payload["parameterType"] = "string"

    with artifact_resource(auth_headers, create_payload) as artifact_id:
        patch_endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
        patch_payload = {
            "artifactType": artifact_type,
            "description": f"Updated description for {artifact_type} ({target_state})",
            "state": target_state,
        }

        # Add type-specific update properties if needed
        if artifact_type == "metric":
            patch_payload["value"] = 0.99  # Updated metric value
        elif artifact_type == "parameter":
            patch_payload["value"] = "0.001"  # Updated parameter value

        patch_response = requests.patch(patch_endpoint, headers=auth_headers, json=patch_payload, timeout=DEFAULT_API_TIMEOUT)
        assert patch_response.status_code == 200
        patch_response_json = patch_response.json()
        assert patch_response_json.get("id") == artifact_id
        assert patch_response_json.get("description") == patch_payload["description"]
        assert patch_response_json.get("state") == patch_payload["state"]
        assert patch_response_json.get("artifactType") == artifact_type
