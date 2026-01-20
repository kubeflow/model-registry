"""Utility functions for search tests."""


def validate_model_contains_search_term(model: dict, search_term: str) -> bool:
    """Verify model contains search term in searchable fields.

    Searchable fields: name, description, provider, libraryName, tasks

    Args:
        model: Model dictionary from API response
        search_term: Search term to validate

    Returns:
        True if model contains search term in any searchable field
    """
    search_term_lower = search_term.lower()

    searchable_content = [
        model.get("name", "").lower(),
        model.get("description", "").lower(),
        model.get("provider", "").lower(),
        model.get("libraryName", "").lower(),
        " ".join(model.get("tasks", [])).lower() if model.get("tasks") else "",
    ]

    return any(search_term_lower in content for content in searchable_content if content)


def get_expected_model_names_for_search(catalog_data: dict, search_term: str) -> set[str]:
    """Get model names that should match the search term from catalog data.

    Args:
        catalog_data: Dictionary containing catalog YAML data with 'models' key
        search_term: Search term to filter models by

    Returns:
        Set of model names that should match the search term
    """
    expected_model_names = set()
    search_lower = search_term.lower()

    for model in catalog_data.get("models", []):
        # Check searchable fields: name, description, provider, libraryName, tasks
        searchable_content = [
            model.get("name", "").lower(),
            model.get("description", "").lower(),
            model.get("provider", "").lower(),
            model.get("libraryName", "").lower(),
            " ".join(model.get("tasks", [])).lower() if model.get("tasks") else "",
        ]

        if any(search_lower in content for content in searchable_content if content):
            expected_model_names.add(model["name"])

    return expected_model_names


def validate_search_results_against_test_data(
    api_response: dict,
    search_term: str,
    catalog_data: dict,
) -> tuple[bool, list[str]]:
    """Validate API search results against expected results from test catalog data.

    Args:
        api_response: API response from search query
        search_term: Search term used
        catalog_data: Test catalog YAML data

    Returns:
        Tuple of (is_valid, list_of_error_messages)
    """
    # Get expected model names from test data
    expected_model_names = get_expected_model_names_for_search(catalog_data, search_term)

    # Get actual model names from API response
    actual_models = api_response.get("items", [])
    actual_model_names = {model["name"] for model in actual_models}

    errors = []

    # Check for missing models
    missing = expected_model_names - actual_model_names
    if missing:
        errors.append(f"Missing models: {sorted(missing)}")

    # Check for extra models
    extra = actual_model_names - expected_model_names
    if extra:
        errors.append(f"Extra models: {sorted(extra)}")

    is_valid = len(errors) == 0
    return is_valid, errors
