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
def generated_schema(request: pytest.FixtureRequest, pytestconfig: pytest.Config,
                     verify_ssl: bool) -> BaseOpenAPISchema:
    """Generate schema for the API based on the schema_file parameter"""
    schema_file = getattr(request, "param", "model-registry.yaml")
    os.environ["API_HOST"] = REGISTRY_URL

    # Read and modify schemathesis.toml if verify_ssl is False
    toml_path = f"{pytestconfig.rootpath}/schemathesis.toml"
    config = schemathesis.config.SchemathesisConfig.from_path(toml_path)
    # tls-verify is by default true
    if verify_ssl is False:
        with open(toml_path) as f:
            toml_content = f.read()

        # Replace tls-verify = true with tls-verify = false
        modified_content = toml_content.replace("tls-verify = true", "tls-verify = false")

        # Write to temporary file
        import tempfile
        with tempfile.NamedTemporaryFile(mode="w", suffix=".toml", delete=False) as temp_file:
            temp_file.write(modified_content)
            temp_toml_path = temp_file.name

        config = schemathesis.config.SchemathesisConfig.from_path(temp_toml_path)

        # Clean up temp file later
        os.unlink(temp_toml_path)
    print(f"Generating schema for {config}")
    schema = schemathesis.openapi.from_path(
        path=f"{pytestconfig.rootpath}/../../api/openapi/{schema_file}",
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
def state_machine(generated_schema: BaseOpenAPISchema, auth_headers: str, pytestconfig: pytest.Config,
                  verify_ssl: bool) -> APIStateMachine:
    BaseAPIWorkflow = generated_schema.as_state_machine()

    class APIWorkflow(BaseAPIWorkflow):  # type: ignore
        headers: dict[str, str]
        verify: bool

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
            self.verify = verify_ssl

        def before_call(self, case: Case) -> None:
            print(f"Checking: {case.method} {case.path}")

        def get_call_kwargs(self, case: Case) -> dict[str, Any]:
            return {"verify": self.verify, "headers": self.headers}

        def after_call(self, response: Response, case: Case) -> None:
            print(f"{case.method} {case.path} -> {response.status_code},")

    return APIWorkflow


@pytest.fixture
def cleanup_artifacts(request: pytest.FixtureRequest, auth_headers: dict, verify_ssl: bool):
    """Cleanup artifacts created during the test."""
    created_ids = []

    def register(artifact_id):
        created_ids.append(artifact_id)

    yield register

    for artifact_id in created_ids:
        del_url = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
        try:
            requests.delete(del_url, headers=auth_headers, timeout=DEFAULT_API_TIMEOUT, verify=verify_ssl)
        except Exception as e:
            print(f"Failed to delete artifact {artifact_id}: {e}")


@pytest.fixture
def artifact_resource(verify_ssl: bool):
    """Create an artifact resource for the test."""

    @contextlib.contextmanager
    def _artifact_resource(auth_headers: dict, payload: dict) -> Generator[str, None, None]:
        create_endpoint = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts"
        resp = requests.post(create_endpoint, headers=auth_headers, json=payload, timeout=DEFAULT_API_TIMEOUT,
                             verify=verify_ssl)
        resp.raise_for_status()
        artifact_id = resp.json()["id"]
        try:
            yield artifact_id
        finally:
            del_url = f"{REGISTRY_URL}/api/model_registry/v1alpha3/artifacts/{artifact_id}"
            try:
                requests.delete(del_url, headers=auth_headers, timeout=DEFAULT_API_TIMEOUT, verify=verify_ssl)
            except Exception as e:
                print(f"Failed to delete artifact {artifact_id}: {e}")

    return _artifact_resource

