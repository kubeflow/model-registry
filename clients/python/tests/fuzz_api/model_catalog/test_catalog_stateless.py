import pytest
import schemathesis
from hypothesis import settings

schema = schemathesis.pytest.from_fixture("generated_schema")

@pytest.mark.parametrize("generated_schema", ["catalog.yaml"], indirect=True)
@schema.parametrize()
@settings(
    deadline=None,
)
@pytest.mark.fuzz
def test_catalog_api_stateless(request_headers: dict[str, str], case: schemathesis.Case):
    """Test the Model Catalog API endpoints.

    This test uses schemathesis to generate and validate API requests
    """
    case.call_and_validate(headers=request_headers)
