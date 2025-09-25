import contextlib
import os
import subprocess
from collections.abc import Generator
from typing import Any

import pytest
import requests
import schemathesis
from schemathesis import Case, Response
from schemathesis.generation.stateful.state_machine import APIStateMachine
from schemathesis.specs.openapi.schemas import BaseOpenAPISchema

from tests.constants import DEFAULT_API_TIMEOUT, REGISTRY_URL


@pytest.fixture
def generated_schema(request: pytest.FixtureRequest, pytestconfig: pytest.Config) -> BaseOpenAPISchema:
    """Generate schema for the API based on the schema_file parameter"""
    schema_file = getattr(request, "param", "model-registry.yaml")
    os.environ["API_HOST"] = REGISTRY_URL
    config = schemathesis.config.SchemathesisConfig.from_path(f"{pytestconfig.rootpath}/schemathesis.toml")
    local_schema_path = f"{pytestconfig.rootpath}/../../api/openapi/{schema_file}"
    schema = schemathesis.openapi.from_path(
        path=local_schema_path,
        config=config,
    )
    schema.config.output.sanitization.update(enabled=False)
    return schema

@pytest.fixture
def auth_headers(user_token: str) -> dict[str, str]:
    """Provides authorization headers for API requests."""
    return {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {user_token}",
    }

@pytest.fixture
def state_machine(generated_schema: BaseOpenAPISchema, auth_headers: str, pytestconfig: pytest.Config) -> APIStateMachine:
    BaseAPIWorkflow = generated_schema.as_state_machine()

    class APIWorkflow(BaseAPIWorkflow):  # type: ignore
        headers: dict[str, str]

        def setup(self) -> None:
            print("Cleaning up database")
            root_path = pytestconfig.rootpath.parent.parent
            cleanup_script = root_path / "scripts" / "cleanup.sh"
            subprocess.run(  # noqa: S603
                [str(cleanup_script)],
                capture_output=True,
                check=True
            )
            self.headers = auth_headers

        def before_call(self, case: Case) -> None:
            print(f"Checking: {case.method} {case.path}")
        def get_call_kwargs(self, case: Case) -> dict[str, Any]:
            return {"verify": False, "headers": self.headers}

        def after_call(self, response: Response, case: Case) -> None:
            print(f"{case.method} {case.path} -> {response.status_code},")
    return APIWorkflow


@pytest.fixture
def cleanup_artifacts(request: pytest.FixtureRequest, auth_headers: dict):
    """Cleanup artifacts created during the test."""
    created_ids = []
    def register(artifact_id):
        created_ids.append(artifact_id)

    yield register

    for artifact_id in created_ids:
        del_url = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
        try:
            requests.delete(del_url, headers=auth_headers, timeout=DEFAULT_API_TIMEOUT)
        except Exception as e:
            print(f"Failed to delete artifact {artifact_id}: {e}")

@pytest.fixture
def artifact_resource():
    """Create an artifact resource for the test."""
    @contextlib.contextmanager
    def _artifact_resource(auth_headers: dict, payload: dict) -> Generator[str, None, None]:
        create_endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts"
        resp = requests.post(create_endpoint, headers=auth_headers, json=payload, timeout=DEFAULT_API_TIMEOUT)
        resp.raise_for_status()
        artifact_id = resp.json()["id"]
        try:
            yield artifact_id
        finally:
            del_url = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
            try:
                requests.delete(del_url, headers=auth_headers, timeout=DEFAULT_API_TIMEOUT)
            except Exception as e:
                print(f"Failed to delete artifact {artifact_id}: {e}")
    return _artifact_resource

