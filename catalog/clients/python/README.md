# Model Catalog Python Client

Python client and E2E tests for the Kubeflow Model Catalog.

## Prerequisites

To run the catalog tests, you need the following tools installed:

- **Python** 3.10, 3.11, or 3.12
- **Poetry** - Python dependency management ([install guide](https://python-poetry.org/docs/#installation))
- **Docker** or **Podman** - Container runtime for building images
- **kubectl** - Kubernetes CLI ([install guide](https://kubernetes.io/docs/tasks/tools/))
- **Kind** - Kubernetes in Docker ([install guide](https://kind.sigs.k8s.io/docs/user/quick-start/#installation))

**Note:** Kustomize 5.5.0 is automatically installed to `bin/` when running deployment targets.

## Installation

```bash
# Install dependencies (Poetry creates a virtual environment automatically)
poetry install

# Or, if you prefer to manage your own virtualenv:
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
poetry install

# Generate OpenAPI client (if needed)
make generate
```

## Running Tests

Tests assume a catalog service is already running (locally, in K8s, etc.).

### Quick Start (Local Kind Cluster)

```bash
# Deploy catalog to a local Kind cluster (creates cluster, builds, deploys)
make deploy

# Run E2E tests
make test-e2e

# Cleanup when done
make deploy-cleanup
```

### Deployment Targets

| Target | Description |
|--------|-------------|
| `make deploy` | Full local deployment (Kind + build + deploy) |
| `make deploy-kind` | Create Kind cluster only |
| `make deploy-build` | Build Docker image only |
| `make deploy-load` | Load image into Kind cluster |
| `make deploy-k8s` | Deploy to existing K8s cluster (no Kind, no build) |
| `make deploy-apply` | Apply kustomize manifests (auto-installs kustomize 5.5.0 if needed) |
| `make deploy-forward` | Start port-forward |
| `make deploy-restart` | Rebuild and restart catalog (after code changes) |
| `make deploy-cleanup` | Remove deployment and Kind cluster |

Deployment uses the kustomize overlay at `manifests/kustomize/options/catalog/overlays/e2e/`.

### Test Commands

| Command | Description |
|---------|-------------|
| `make test-e2e` | Run E2E tests against running service |
| `make test-fuzz` | Run fuzz tests against running service |
| `make nox-e2e` | Run E2E tests on all Python versions |

### Using pytest directly

```bash
# Run E2E tests
poetry run pytest --e2e

# Run fuzz tests (API schema fuzzing)
poetry run pytest --fuzz

# Run specific test file
poetry run pytest tests/test_ordering.py --e2e

# Run with verbose output
poetry run pytest --e2e -v -rA

# Run in parallel
poetry run pytest --e2e -n auto
```

## Testing Against External K8s

If you have a catalog service running in K8s (not using make deploy):

```bash
# Port-forward the catalog service
kubectl port-forward svc/catalog 8081:8080 -n model-registry

# Set environment variables and run tests
export CATALOG_URL="http://localhost:8081"
poetry run pytest --e2e
```

## Multi-Python Version Testing with Nox

```bash
# List available nox sessions
poetry run nox -l

# Run default sessions (lint + tests)
make nox

# Run E2E tests on all Python versions
make nox-e2e
```

## Development

```bash
# Lint code
make lint

# Auto-fix code style
make tidy

# Build package
make build

# Update dependencies
make update
```

## Project Structure

```
catalog/clients/python/
├── pyproject.toml          # Poetry configuration & dependencies
├── poetry.lock             # Lock file with pinned versions
├── noxfile.py              # Nox sessions for multi-Python testing
├── Makefile                # Build & deployment automation
├── README.md               # This file
├── schemathesis.toml       # Schemathesis (API fuzzing) config
├── patches/                # Patches for generated code
├── src/
│   ├── model_catalog/      # Client library
│   │   ├── __init__.py     # Package exports
│   │   └── _client.py      # CatalogAPIClient
│   └── catalog_openapi/    # Generated OpenAPI client
│       ├── api/            # API classes
│       └── models/         # Model classes
└── tests/
    ├── conftest.py         # Pytest fixtures & configuration
    ├── constants.py        # Test constants & configuration
    ├── fuzz_api/           # API fuzzing tests (Schemathesis)
    ├── test_artifacts.py   # Artifact filtering & ordering tests
    ├── test_filter_options.py  # Filter options & named queries tests
    ├── test_models.py      # Model filtering tests
    ├── test_ordering.py    # Name & accuracy ordering tests
    ├── test_source_preview.py  # Source preview tests
    └── test_sources.py     # Source status & configuration tests

# K8s manifests are in the main kustomize directory:
manifests/kustomize/options/catalog/overlays/e2e/
├── kustomization.yaml      # Kustomize overlay for E2E testing
├── sources.yaml            # Test catalog sources
├── test-catalog.yaml       # Main test catalog data
└── disabled-catalog.yaml   # Disabled source catalog data
```

## Nox Sessions

Nox runs tests on Python 3.10, 3.11, and 3.12.

| Session | Description |
|---------|-------------|
| `lint` | Lint using ruff |
| `mypy` | Type check using mypy |
| `e2e` | Run E2E tests |
| `coverage` | Produce coverage report |

## Test Markers

Tests are marked with pytest markers:

| Marker | Description |
|--------|-------------|
| `@pytest.mark.e2e` | End-to-end tests (require running service) |
| `@pytest.mark.fuzz` | API fuzzing tests (require running service) |
| `@pytest.mark.huggingface` | Tests that interact with HuggingFace API |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CATALOG_URL` | Catalog service URL | `http://localhost:8081` |
| `CATALOG_NAMESPACE` | K8s namespace | `model-catalog` |
| `CATALOG_IMAGE` | Docker image name | `model-registry:catalog-test` |
| `CLUSTER_NAME` | Kind cluster name | `catalog-e2e` |
| `CATALOG_PORT` | Local port for port-forward | `8081` |
