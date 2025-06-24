import pytest
import schemathesis
from hypothesis import HealthCheck, settings

schema = schemathesis.pytest.from_fixture("generated_schema")


@schema.parametrize()
@settings(
    deadline=None,
    suppress_health_check=[
        HealthCheck.filter_too_much,
        HealthCheck.too_slow,
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

    case.call_and_validate(headers=headers)

