import pytest
import schemathesis
from hypothesis import HealthCheck, settings

schema = schemathesis.pytest.from_fixture("generated_schema")


@schema.parametrize()
@settings(
    deadline=None,
    suppress_health_check=[
        HealthCheck.filter_too_much,
    ]
)
@pytest.mark.e2e
def test_mr_api_stateless(setup_env_user_token, case):
    """Test the Model Registry API endpoints.

    This test uses schemathesis to generate and validate API requests
    """

    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {setup_env_user_token}"
    }

    # We need to ensure valid data
    if case.path == "/api/model_registry/v1alpha3/model_versions/{modelversionId}/artifacts":
        case.path_parameters["modelversionId"] = "test-model-version-1"

        if case.method == "POST":
            case.body = {
                "artifactType": "model-artifact",
                "name": "test-artifact",
                "uri": "s3://test-bucket/test-artifact",
                "metadata": {}
            }

    case.call_and_validate(headers=headers)

