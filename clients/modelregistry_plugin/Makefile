# Makefile for Model Registry MLflow Plugin

.PHONY: help install build test test-e2e test-e2e-local clean lint format

# Default target
help:
	@echo "Available targets:"
	@echo "  install      - Install dependencies"
	@echo "  build        - Build the package"
	@echo "  install-dev  - Install in development mode"
	@echo "  test         - Run unit tests with coverage"
	@echo "  test-e2e     - Run end-to-end tests (requires MLFLOW_TRACKING_URI)"
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
	uv sync --dev

# Run unit tests
test:
	uv run pytest tests/ -v --tb=short --cov=modelregistry_plugin --cov-report=term-missing -k "not e2e"

# Run end-to-end tests (requires MLFLOW_TRACKING_URI)
test-e2e:
	@echo "Running end-to-end tests..."
	@if [ -z "$$MLFLOW_TRACKING_URI" ]; then \
		echo "❌ MLFLOW_TRACKING_URI is not set"; \
		echo "Please set: export MLFLOW_TRACKING_URI='modelregistry://host:port'"; \
		exit 1; \
	fi
	@echo "✅ Using tracking URI: $$MLFLOW_TRACKING_URI"
	@if [ -n "$$MODEL_REGISTRY_TOKEN" ]; then \
		echo "✅ Using authentication token"; \
	else \
		echo "⚠️  No MODEL_REGISTRY_TOKEN set (may fail if auth required)"; \
	fi
	uv run pytest tests/test_e2e.py -v -s

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
	uv run python -c "import mlflow; print(list(mlflow.tracking._tracking_service.utils._tracking_store_registry._registry.keys()))"

# Verify entry point registration
verify-entry-point:
	@echo "Verifying entry point registration..."
	uv run python -c "import mlflow; stores = list(mlflow.tracking._tracking_service.utils._tracking_store_registry._registry.keys()); print('modelregistry' in stores and '✅ Entry point registered' or '❌ Entry point not found')"

# Full development workflow
dev: install build install-dev verify-entry-point test 
