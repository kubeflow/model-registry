# MLflow Model Registry Store Plugin

## Quick demo

```
git clone -b feat/mlflow-model-registry-store --single-branch https://github.com/jonburdo/model-registry.git
cd clients/mlflow-model-registry-store

./demo-native.sh
./demo-modelregistry.sh
```

For `demo-modelregistry.sh` you will need a local model registry instance running. From repo root:
```
cd clients/python
make deploy-latest-mr deploy-local-registry deploy-test-minio
```

This package provides an MLflow plugin that allows using a Model Registry server as the backend for MLflow's model registry operations.

## Installation

With repo url:
```bash
pip install git+https://github.com/jonburdo/model-registry.git@feat/mlflow-model-registry-store#subdirectory=clients/mlflow-model-registry-store
```

Or clone this repo and then install in install:

```bash
git clone -b feat/mlflow-model-registry-store --single-branch https://github.com/jonburdo/model-registry.git
cd clients/mlflow-model-registry-store

cd clients/mlflow-model-registry-store
uv pip install .
```
or use `-e` to install in development mode:
```
uv pip install -e .
```

## Usage

Once installed, you can configure MLflow to use the Model Registry as follows:

### URI Formats

The plugin supports multiple URI formats for flexibility:

```python
import os
import mlflow

# Option 1: Basic scheme (defaults to HTTPS)
os.environ["MLFLOW_REGISTRY_URI"] = "modelregistry://my-server:8080?author=myuser"

# Option 2: Explicit HTTP
os.environ["MLFLOW_REGISTRY_URI"] = "modelregistry+http://my-server:8080?author=myuser"

# Option 3: Explicit HTTPS
os.environ["MLFLOW_REGISTRY_URI"] = "modelregistry+https://my-server:8443?author=myuser"

# Option 4: Override security with query parameter
os.environ["MLFLOW_REGISTRY_URI"] = "modelregistry://my-server:8080?author=myuser&is-secure=false"
```

### Environment Variable Configuration

For easier deployment and configuration management, you can use environment variables as defaults:

```bash
# Environment variables (used as fallbacks)
export MODEL_REGISTRY_HOST=my-server
export MODEL_REGISTRY_PORT=8080
export MODEL_REGISTRY_AUTHOR=myuser
export MODEL_REGISTRY_SECURE=false
export MODEL_REGISTRY_TOKEN=your-token

# Minimal URI when using environment variables
export MLFLOW_REGISTRY_URI="modelregistry:///"
```

**Priority order**: URL parameters > Environment variables > Defaults

### Query Parameters

| Parameter | Environment Variable | Description | Default |
|-----------|----------------------|-------------|---------|
| `author` | `MODEL_REGISTRY_AUTHOR` | Author name for model operations | `"unknown"` |
| `is-secure` | `MODEL_REGISTRY_SECURE` | Use HTTPS instead of HTTP | `true` |
| `user-token` | `MODEL_REGISTRY_TOKEN` | Authentication token | `None` |
| `custom-ca` | `MODEL_REGISTRY_CA` | Custom CA certificate path | `None` |

```python
# Now use MLflow model registry as usual
with mlflow.start_run():
    # Log and register a model
    mlflow.sklearn.log_model(
        model,
        "model",
        registered_model_name="my-model"
    )

# Get registered model
client = mlflow.MlflowClient()
model = client.get_registered_model("my-model")
```

## URI Format

The store URI follows this format:

```
modelregistry://server:port?author=<author>&is-secure=<true/false>&user-token=<token>&custom-ca=<ca_path>
```

### Parameters:

- **server**: Model Registry server hostname
- **port**: Server port (default: 443)
- **author**: Author name for model registration (required)
- **is-secure**: Whether to use HTTPS (default: true)
- **user-token**: Optional authentication token
- **custom-ca**: Optional path to custom CA certificate

### Examples:

```bash
# Basic configuration
export MLFLOW_REGISTRY_URI="modelregistry://localhost:8080?author=myuser&is-secure=false"

# With authentication
export MLFLOW_REGISTRY_URI="modelregistry://mr-server.model-registry.svc.cluster.local:443?author=myuser&user-token=/var/run/secrets/kubernetes.io/serviceaccount/token"
```

## Features

This plugin implements the MLflow AbstractStore interface and provides:

-  Create/get/search registered models
-  Create/get/update model versions
-  Set/delete tags on models and versions
-  Model version staging via tags
-  Model aliases via custom properties
- ✅ Delete operations (implemented via Model Registry archive)
- ❌ Model renaming (not supported by Model Registry)

## Mapping Between MLflow and Model Registry

| MLflow Concept | Model Registry Concept | Notes |
|----------------|-------------------|-------|
| RegisteredModel | RegisteredModel | Direct mapping |
| ModelVersion | ModelVersion + ModelArtifact | Combined for source URI |
| Tags | CustomProperties | Stored as key-value pairs |
| Stages | CustomProperties | Stored as `mlflow.stage` property |
| Aliases | CustomProperties | Stored as `mlflow.alias.<name>` |
| Run ID | CustomProperties | Stored as `mlflow.run_id` |

## Delete Operations

MLflow delete operations are implemented using Model Registry's archive functionality:

- **`delete_registered_model(name)`**: Sets the model state to `ARCHIVED` and archives all its versions
- **`delete_model_version(name, version)`**: Sets the specific version state to `ARCHIVED`

Archived models and versions:
- Are excluded from search results and listings
- Cannot be retrieved via `get_registered_model()` or `get_model_version()`
- Are preserved in the Model Registry database (soft delete)
- Can be restored by updating the state back to `LIVE` using the Model Registry client directly

This approach ensures MLflow delete semantics while maintaining Model Registry's immutable design principles.

## Development

```bash
# Install development dependencies
uv pip install -e ".[dev]"

# Run unit tests (with mocked backend)
pytest tests/test_store.py

# Run integration tests (requires running Model Registry server)
# First, start Model Registry server at 127.0.0.1:8080
kubectl port-forward svc/model-registry-service 8080:8080

# Then run integration tests
uv run pytest tests/test_integration.py -v

# Run MLflow model registry tests against this backend
export MLFLOW_REGISTRY_URI="modelregistry://localhost:8080?author=testuser&is-secure=false"
python -m pytest path/to/mlflow/tests/store/model_registry/
```

### Integration Tests

The integration tests in `tests/test_integration.py` connect to a real Model Registry server without any mocking. To run these tests:

1. **Start Model Registry Server**: Ensure you have a Model Registry server running at `http://127.0.0.1:8080`
   ```bash
   # Using kubectl port-forward (if running in Kubernetes)
   kubectl port-forward svc/model-registry-service 8080:8080

   # Or start locally if you have the server binary
   # (follow Model Registry server documentation)
   ```

2. **Run Integration Tests**:
   ```bash
   uv run pytest tests/test_integration.py -v
   ```

3. **Test Individual Features**:
   ```bash
   # Test just connection
   uv run pytest tests/test_integration.py::TestModelRegistryIntegration::test_connection_and_client_creation -v

   # Test model registration
   uv run pytest tests/test_integration.py::TestModelRegistryIntegration::test_model_registration_and_retrieval -v
   ```

**Note**: Integration tests will create and modify data in the connected Model Registry instance. The tests use unique model names to avoid conflicts but do not clean up models (as Model Registry doesn't support deletion).

## Limitations

1. **No Model Renaming**: Model names are immutable in Model Registry
2. **Limited Search**: Cross-model version search is not optimized
3. **Stage Implementation**: Stages are implemented via tags, not natively supported
4. **Version Numbering**: Model Registry uses version names, not sequential numbers
5. **Soft Delete Only**: Delete operations archive rather than permanently remove data

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request
