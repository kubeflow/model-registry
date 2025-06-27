import pytest
import schemathesis
from hypothesis import HealthCheck, settings, strategies as st
from schemathesis.specs.openapi.formats import register_string_format

schema = schemathesis.pytest.from_fixture("generated_schema")


int64_string_strategy_instance = st.integers(min_value=1, max_value=2**31 - 1).map(str)

register_string_format("int64", int64_string_strategy_instance)

@schema.parametrize()
@settings(
    deadline=None,
    suppress_health_check=[
        HealthCheck.filter_too_much,
        HealthCheck.too_slow,
    ],
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
