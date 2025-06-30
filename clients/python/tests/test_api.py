import pytest
import schemathesis
from hypothesis import HealthCheck, settings, Verbosity
from hypothesis import strategies as st
from schemathesis.specs.openapi.formats import register_string_format
from schemathesis.hooks import hook

schema = schemathesis.pytest.from_fixture("generated_schema")

# Strategy for int64 string format
int64_string_strategy_instance = st.integers(min_value=1, max_value=2**63 - 1).map(str)
register_string_format("int64", int64_string_strategy_instance)

# Custom strategies for artifact fields
artifact_type_choices = st.sampled_from(["model-artifact", "doc-artifact"])
artifact_state_choices = st.sampled_from(["LIVE", "PENDING", "MARKED_FOR_DELETION", "DELETED", "ABANDONED", "REFERENCE", "UNKNOWN"])

@hook
def before_generate_case(context, strategy):
    """
    Hook to modify the strategy before test case generation
    Unsatisfiable error with the default strategy for the following endpoints:
    - PATCH /api/v1/artifacts/{id}
    - POST /api/v1/model_versions/{modelversionId}/artifacts
    """
    if context.operation.method == "patch" and "artifacts/{id}" in context.operation.path:
        print(f"Applying custom strategy for PATCH artifacts endpoint for {context.operation.path}")

        def modify_case(case):
            print(f"\nModifying case for PATCH:")
            print(f"Before modification: {case}")

            # Always ensure we have a valid body structure
            if not case.body:
                case.body = {}

            # For positive tests, always provide valid required fields
            case.body = {
                "artifactType": "model-artifact",
                "state": "LIVE",
                "uri": "s3://test-bucket/model.pkl",
                "description": "Test artifact"
            }

            # If there are custom properties, ensure they have valid structure
            if "customProperties" in case.body:
                case.body["customProperties"] = {
                    "test_string": {
                        "metadataType": "MetadataStringValue",
                        "string_value": "test"
                    }
                }

            print(f"After modification: {case}")
            return case

        modified_strategy = strategy.map(modify_case)
        print(f"Strategy after map: {modified_strategy}")
        return modified_strategy

    elif context.operation.method == "post" and "model_versions/{modelversionId}/artifacts" in context.operation.path:
        print(f"Applying custom strategy for POST model_versions artifacts endpoint: {context.operation.path}")

        # For POST, we need more fields including artifactType
        def modify_case(case):
            print(f"\nModifying case for POST:")
            print(f"Before modification: {case}")

            # Always ensure we have a valid body structure
            if not case.body:
                case.body = {}

            # Provide valid required fields for POST
            case.body = {
                "artifactType": "model-artifact",
                "name": "test-artifact",
                "state": "LIVE",
                "uri": "s3://test-bucket/model.pkl",
                "description": "Test artifact"
            }

            # If there are custom properties, ensure they have valid structure
            if "customProperties" in case.body:
                case.body["customProperties"] = {
                    "test_string": {
                        "metadataType": "MetadataStringValue",
                        "string_value": "test"
                    }
                }

            print(f"After modification: {case}")
            return case

        modified_strategy = strategy.map(modify_case)
        print(f"Strategy after map: {modified_strategy}")
        return modified_strategy

    return strategy

@schema.parametrize()
@settings(
    deadline=None,
    suppress_health_check=[
        HealthCheck.filter_too_much,
        HealthCheck.too_slow,
        HealthCheck.data_too_large,
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
