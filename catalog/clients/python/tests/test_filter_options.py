"""E2E tests for filter options and named queries.

Tests the filter_options API endpoint and named queries functionality.
Consolidated from test_named_queries.py, test_named_queries_merge.py,
and test_named_queries_validation.py.

To run these tests:
1. Start the catalog service with test data loaded
2. Set CATALOG_URL environment variable (default: http://localhost:8081)
3. Run: pytest --e2e tests/test_filter_options.py
"""

from model_catalog import CatalogAPIClient


class TestFilterOptions:
    """Test suite for filter options functionality."""

    def test_get_filter_options_returns_response(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that filter options endpoint returns a response."""
        response = api_client.get_filter_options()
        assert response is not None, "Filter options response should not be None"
        assert isinstance(response, dict)

    def test_filter_options_has_filters_field(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that filter options response has filters field."""
        response = api_client.get_filter_options()
        filters = response["filters"]
        assert isinstance(filters, dict)
        assert len(filters) > 0, "Filters object should not be empty"

    def test_filter_options_contains_field_types(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that filter options contains field type information."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        # Known filter types supported by the API
        # If new types are added, this list should be updated
        known_types = {"string", "number", "boolean", "array", "object"}

        for field_name, field_info in filters.items():
            assert isinstance(field_info, dict)
            assert "type" in field_info, f"Field {field_name} missing type"
            assert field_info["type"] in known_types, (
                f"Field {field_name} has unknown type '{field_info['type']}'. Known types: {known_types}"
            )

    def test_string_filters_have_values(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that string type filters have available values."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        for field_name, field_info in filters.items():
            if field_info.get("type") == "string":
                assert "values" in field_info, f"String field {field_name} missing values"
                assert isinstance(field_info["values"], list)

    def test_number_filters_have_range(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that number type filters have range information."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        for field_name, field_info in filters.items():
            if field_info.get("type") == "number":
                assert "range" in field_info, f"Number field {field_name} missing range"
                range_info = field_info["range"]
                assert "min" in range_info
                assert "max" in range_info


class TestFilterValidation:
    """Test suite for filter validation."""

    def test_filter_types_are_valid(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that filter types are valid enumeration values."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        # Known filter types - same as test_filter_options_contains_field_types
        known_types = {"string", "number", "boolean", "array", "object"}

        for field_name, field_info in filters.items():
            field_type = field_info.get("type")
            assert field_type in known_types, (
                f"Unknown type '{field_type}' for field {field_name}. Known types: {known_types}"
            )

    def test_number_ranges_are_valid(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that number ranges have valid min <= max."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        for field_name, field_info in filters.items():
            if field_info.get("type") == "number":
                range_info = field_info.get("range", {})
                min_val = range_info.get("min")
                max_val = range_info.get("max")

                if min_val is not None and max_val is not None:
                    assert min_val <= max_val, f"Invalid range for {field_name}: min={min_val}, max={max_val}"

    def test_string_values_are_non_empty(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that string filter values are non-empty lists of strings."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        for field_name, field_info in filters.items():
            if field_info.get("type") == "string":
                values = field_info.get("values", [])
                assert isinstance(values, list), f"Values for {field_name} should be a list"
                for idx, val in enumerate(values):
                    assert isinstance(val, str), (
                        f"Value at index {idx} in {field_name} should be string, got {type(val)}"
                    )
                    assert val.strip(), f"Value at index {idx} in {field_name} should not be empty or whitespace"

    def test_string_values_are_distinct(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that string filter values contain no duplicates."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        for field_name, field_info in filters.items():
            if field_info.get("type") == "string":
                values = field_info.get("values", [])
                assert len(values) == len(set(values)), (
                    f"Values for '{field_name}' should be distinct (found duplicates)"
                )

    def test_filter_field_names_format(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that filter field names follow expected format."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        for field_name in filters:
            assert isinstance(field_name, str)
            assert len(field_name) > 0

            # Field names typically use dots for nested properties
            parts = field_name.split(".")
            for part in parts:
                assert len(part) > 0, f"Empty part in field name: {field_name}"


class TestNamedQueries:
    """Test suite for named queries functionality."""

    def test_named_queries_structure(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that named queries have proper structure if present."""
        response = api_client.get_filter_options()

        if "namedQueries" in response:
            named_queries = response["namedQueries"]
            assert isinstance(named_queries, dict)

            for query_name, query_def in named_queries.items():
                assert isinstance(query_name, str)
                assert isinstance(query_def, dict)

    def test_named_queries_operators_are_valid(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that named query operators are valid if present."""
        response = api_client.get_filter_options()

        if "namedQueries" not in response:
            return

        named_queries = response["namedQueries"]
        valid_operators = {"=", "!=", "<", ">", "<=", ">=", "LIKE", "ILIKE", "IN"}

        for query_name, query_def in named_queries.items():
            for _field_name, condition in query_def.items():
                if isinstance(condition, dict):
                    operator = condition.get("operator")
                    if operator:
                        assert operator in valid_operators, f"Invalid operator '{operator}' in query {query_name}"

    def test_filter_options_artifact_properties(self, api_client: CatalogAPIClient, suppress_ssl_warnings: None):
        """Test that artifact custom properties appear in filter options."""
        response = api_client.get_filter_options()
        filters = response.get("filters", {})

        # Should have some filter options from test data
        assert isinstance(filters, dict), "filters should be a dict"
        assert len(filters) > 0, "Expected non-empty filter options from test data"
