# Makefile for Model Registry MLflow Plugin

.PHONY: help install build test test-e2e test-e2e-setup clean lint format test-e2e-local

# Default target
help:
	@echo "Available targets:"
	@echo "  install      - Install dependencies"
	@echo "  build        - Build the package"
	@echo "  install-dev  - Install in development mode"
	@echo "  test         - Run unit tests"
	@echo "  test-e2e     - Run end-to-end tests"
	@echo "  test-e2e-setup - Test e2e setup and connectivity"
	@echo "  test-e2e-local - Run local end-to-end tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  lint         - Run linting"
	@echo "  format       - Format code"

# Install dependencies
install:
	uv sync

# Build the package
build:
	uv build

# Install in development mode
install-dev:
	uv pip install -e .

# Run unit tests
test:
	uv run pytest tests/ -v --tb=short

# Run end-to-end tests
test-e2e:
	@echo "Running end-to-end tests..."
	@echo "Make sure you have set the required environment variables:"
	@echo "  MODEL_REGISTRY_HOST"
	@echo "  MODEL_REGISTRY_TOKEN"
	@echo ""
	@echo "You can also use: ./scripts/run_e2e_tests.sh"
	uv run pytest tests/test_e2e.py -v -s

# Test e2e setup and connectivity
test-e2e-setup:
	@echo "Testing e2e setup and connectivity..."
	python scripts/test_e2e_setup.py

# Run local end-to-end tests (starts local MLMD and Model Registry servers)
test-e2e-local:
	@echo "Running local end-to-end tests (starts local MLMD and Model Registry servers)..."
	uv run pytest tests/test_e2e_local.py -v -s

# Clean build artifacts
clean:
	rm -rf dist/
	rm -rf build/
	rm -rf *.egg-info/
	find . -type d -name __pycache__ -delete
	find . -type f -name "*.pyc" -delete

# Run linting
lint:
	uv run ruff check .

# Format code
format:
	uv run ruff format .

# Show available tracking stores
show-stores:
	@echo "Available MLflow tracking stores:"
	python -c "import mlflow; print(list(mlflow.tracking._tracking_service.utils._tracking_store_registry._registry.keys()))"

# Verify entry point registration
verify-entry-point:
	@echo "Verifying entry point registration..."
	python -c "import mlflow; stores = list(mlflow.tracking._tracking_service.utils._tracking_store_registry._registry.keys()); print('modelregistry' in stores and '✅ Entry point registered' or '❌ Entry point not found')"

# Full development workflow
dev: install build install-dev verify-entry-point test 