"""E2E tests for the Model Catalog Python client.

Test Structure:
    - test_artifacts.py: Artifact filtering and ordering tests
    - test_filter_options.py: Filter options and named queries tests
    - test_models.py: Model listing and filtering tests
    - test_ordering.py: Model ordering (NAME, ACCURACY) tests
    - test_sources.py: Source management and status tests
    - test_source_preview.py: Source preview functionality tests
    - fuzz_api/: API fuzzing tests using Schemathesis

Running Tests:
    make test-e2e       # Run all E2E tests
    make test-fuzz      # Run fuzzing tests

Requirements:
    - Catalog service running (use `make deploy` for local K8s)
    - CATALOG_URL environment variable (default: http://localhost:8081)
"""

__all__: list[str] = []
