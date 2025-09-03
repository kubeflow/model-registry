import pytest
from hypothesis import HealthCheck, settings


@pytest.mark.fuzz
class TestRestAPIStateful:
    @pytest.mark.parametrize("generated_schema", ["model-registry.yaml"], indirect=True)
    def test_mr_api_stateful(self, state_machine):
        """Launches stateful tests against the Model Registry API endpoints defined in its openAPI yaml spec file"""
        state_machine.run(settings=settings(
            max_examples=20,
            deadline=10000, #10 seconds
            suppress_health_check=[
                HealthCheck.filter_too_much,
                HealthCheck.too_slow,
                HealthCheck.data_too_large,
            ],
        ))
