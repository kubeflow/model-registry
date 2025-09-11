# Kubeflow Model Registry MLflow Plugin

An MLflow tracking plugin that integrates with Kubeflow Model Registry, enabling seamless experiment tracking and model management through MLflow's familiar API while storing metadata in the Kubeflow Model Registry backend.

## Overview

This plugin allows you to use MLflow's tracking API to log experiments, runs, metrics, parameters, and artifacts while having all metadata stored in Kubeflow Model Registry instead of MLflow's default backend storage.

## Features

- **Seamless Integration**: Use standard MLflow tracking APIs with Model Registry backend
- **Experiment Management**: Create and manage experiments in Model Registry
- **Run Tracking**: Log runs, metrics, parameters, and tags
- **Artifact Storage**: Store and retrieve model artifacts
- **Authentication**: Built-in support for Kubeflow authentication
- **SSL/TLS Support**: Automatic CA certificate detection for secure HTTPS connections
- **Kubernetes Ready**: Auto-detects Kubernetes CA certificates and service account tokens
- **Full MLflow Compatibility**: Implements all required AbstractStore methods

## Installation

### Using uv (Recommended)

```bash
# Build from source
cd clients/model-registry-mlflow
uv build
uv pip install dist/model_registry_mlflow-0.1.0-py3-none-any.whl
```

### Using pip

```bash
pip install model-registry-mlflow
```

## Configuration

### Environment Variables

Set the following environment variables to configure the plugin:

```bash
# Required: MLflow tracking URI for the Model Registry server
export MLFLOW_TRACKING_URI="modelregistry://your-model-registry-host:port"

# Optional: Authentication token (if required)
export MODEL_REGISTRY_TOKEN="your-token"
```

### Authentication

The plugin supports various authentication methods:

```bash
# Token-based auth
export MODEL_REGISTRY_TOKEN="your-token"

# Kubernetes service account token from `/var/run/secrets/kubernetes.io/serviceaccount/token`
```

### TLS and CA Certificate Configuration

For secure HTTPS connections to Model Registry, the plugin supports custom CA certificate configuration. This is particularly useful when using self-signed certificates or custom certificate authorities.

#### CA Certificate Priority Order

The plugin automatically detects and configures CA certificates in the following priority order:

1. **Custom CA via Environment Variable** (highest priority)
2. **Kubernetes Default CA** (auto-detected when running in K8s)
3. **System Default CA Bundle** (fallback)

#### Configuration Options

```bash
# Option 1: Custom CA certificate via environment variable
export MODELREGISTRY_CA_CERT_PATH="/path/to/your/ca.crt"
export MLFLOW_TRACKING_URI="modelregistry+https://registry.example.com:8080"

# Option 2: In Kubernetes - works automatically (no configuration needed)
export MLFLOW_TRACKING_URI="modelregistry+https://model-registry-service:8080"
# Plugin automatically uses /run/secrets/kubernetes.io/serviceaccount/ca.crt

# Option 3: HTTP (no TLS) - CA configuration is skipped
export MLFLOW_TRACKING_URI="modelregistry://registry.example.com:8080"
```

#### Kubernetes Deployment Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mlflow-app
spec:
  template:
    spec:
      containers:
      - name: mlflow-app
        image: your-mlflow-app:latest
        env:
        - name: MLFLOW_TRACKING_URI
          value: "modelregistry+https://model-registry-service:8080"
        # CA certificate is automatically detected from Kubernetes service account
```

#### Custom CA in Kubernetes

For custom CA certificates in Kubernetes deployments:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ca-certificates
data:
  ca.crt: |
    -----BEGIN CERTIFICATE-----
    ... your CA certificate content ...
    -----END CERTIFICATE-----
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mlflow-app
spec:
  template:
    spec:
      containers:
      - name: mlflow-app
        image: your-mlflow-app:latest
        env:
        - name: MLFLOW_TRACKING_URI
          value: "modelregistry+https://model-registry-service:8080"
        - name: MODELREGISTRY_CA_CERT_PATH
          value: "/etc/ssl/certs/ca.crt"
        volumeMounts:
        - name: ca-certs
          mountPath: /etc/ssl/certs
          readOnly: true
      volumes:
      - name: ca-certs
        configMap:
          name: ca-certificates
```

#### Protocol Support

- **HTTP** (`modelregistry://`): No CA certificate configuration
- **HTTPS** (`modelregistry+https://`): Automatic CA certificate detection and configuration

## Usage

### Basic Example

```python
import mlflow

# Set tracking URI to use the Model Registry plugin
mlflow.set_tracking_uri("modelregistry://localhost:8080")

# Create experiment
experiment_id = mlflow.create_experiment("my-experiment")

# Start a run
with mlflow.start_run():
    # Log parameters
    mlflow.log_param("learning_rate", 0.01)
    mlflow.log_param("epochs", 100)
    
    # Log metrics
    mlflow.log_metric("accuracy", 0.95)
    mlflow.log_metric("loss", 0.05)
    
    # Log model
    mlflow.sklearn.log_model(model, "model")
```

### Advanced Configuration

```python
import mlflow
from model_registry_mlflow import ModelRegistryStore

# Configure with custom settings
store = ModelRegistryStore(
    store_uri="modelregistry://model-registry.kubeflow.svc.cluster.local:8080",
    artifact_uri="s3://my-bucket/artifacts"
)

# The store is automatically registered with MLflow
mlflow.set_tracking_uri("modelregistry://model-registry.kubeflow.svc.cluster.local:8080")
```

## Plugin Architecture

The plugin implements MLflow's tracking store interface by:

1. **Store Implementation**: `ModelRegistryStore` class implements `AbstractStore`
2. **API Translation**: Converts MLflow API calls to Model Registry API requests
3. **Authentication**: Handles Kubeflow authentication mechanisms
4. **Artifact Management**: Manages artifact storage and retrieval
5. **Circular Import Prevention**: Uses lazy imports to avoid circular dependencies during MLflow initialization

## Supported Operations

### Experiment Operations
- ✅ Create experiments
- ✅ Get experiments by ID or name
- ✅ List experiments with filtering (ACTIVE_ONLY, DELETED_ONLY, ALL)
- ✅ Search experiments with pagination
- ✅ Delete/restore experiments
- ✅ Rename experiments

### Run Operations
- ✅ Create runs
- ✅ Get runs by ID
- ✅ Update run information
- ✅ Delete/restore runs
- ✅ Search runs with filtering and pagination

### Metrics and Parameters
- ✅ Log metrics with timestamps and steps
- ✅ Get metric history for specific metrics
- ✅ Log parameters
- ✅ Batch logging of metrics, parameters, and tags

### Artifacts and Models
- ✅ Log inputs (datasets and models)
- ✅ Log outputs (models)
- ✅ Record logged models
- ✅ Create and manage logged models
- ✅ Search logged models
- ✅ Set and delete model tags

## Unsupported Features

### Trace Operations
The following MLflow trace-related methods are **not supported** in the ModelRegistryStore implementation:

- `start_trace()` - Start a trace
- `end_trace()` - End a trace
- `delete_traces()` - Delete traces
- `get_trace_info()` - Get trace information
- `search_traces()` - Search traces
- `set_trace_tag()` - Set trace tags
- `delete_trace_tag()` - Delete trace tags

These methods are not required for basic MLflow tracking functionality and are typically used for advanced inference monitoring and debugging. The ModelRegistryStore focuses on core experiment tracking and model management features.

### Advanced Features
- **Filter String Support**: Some search methods don't fully support MLflow's filter string syntax (marked as TODO in the implementation)
- **Order By Support**: Advanced ordering options are not fully implemented
- **Batch Operations**: Some batch operations may not be optimized for the Model Registry backend

## Development

For developers who want to contribute to the plugin, see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development setup, testing, and contribution guidelines.

### Quick Development Setup

```bash
# Install dependencies
uv sync

# Install in development mode
uv pip install -e .

# Run tests
make test-e2e-local  # Local e2e tests (recommended)
make test            # Unit tests
```

### Testing

The plugin includes comprehensive testing:

- **Unit Tests**: Fast, isolated tests for individual components
- **Local E2E Tests**: Self-contained tests with local Model Registry server
- **Remote E2E Tests**: Tests against real remote servers (optional)

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed testing information and setup instructions.

## Technical Details

### Circular Import Resolution

The plugin uses lazy imports to prevent circular dependencies during MLflow's entry point registration:

- MLflow imports are moved inside methods where they're needed
- Type annotations are simplified to avoid referencing MLflow types at module level
- This ensures the entry point can be registered without conflicts

### Entry Point Registration

The plugin registers itself as an MLflow tracking store via the entry point:

```toml
[project.entry-points."mlflow.tracking_store"]
modelregistry = "model_registry_mlflow.tracking_store:ModelRegistryStore"
```

This allows MLflow to automatically discover and use the plugin when the `modelregistry://` URI scheme is specified.

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure Model Registry is running and accessible
2. **Authentication Failed**: Verify tokens and credentials are correct
3. **SSL/TLS Certificate Errors**: 
   - Check that the CA certificate path is correct and the file exists
   - Verify the certificate is in PEM format
   - Ensure the server certificate is signed by the specified CA
   - For Kubernetes: verify the service account has access to CA certificates
4. **Entry Point Not Found**: Ensure the package is properly installed and the entry point is registered

### Debug Logging

Enable debug logging to troubleshoot issues:

```python
import logging
logging.getLogger("model_registry_mlflow").setLevel(logging.DEBUG)

# For detailed CA certificate configuration logging:
logging.getLogger("model_registry_mlflow.api_client").setLevel(logging.INFO)
```

**CA Certificate Log Messages:**
- `"Using CA certificate from environment variable MODELREGISTRY_CA_CERT_PATH: /path/to/ca.crt"` - Custom CA detected
- `"Using Kubernetes default CA certificate: /run/secrets/kubernetes.io/serviceaccount/ca.crt"` - K8s CA auto-detected  
- `"Using system default CA bundle"` - Fallback to system CA
- `"SSL certificate verification failed connecting to Model Registry"` - Certificate verification error

### Verification Commands

```python
# Check if the plugin is registered
import mlflow
print('Available tracking stores:', list(mlflow.tracking._tracking_service.utils._tracking_store_registry._registry.keys()))

# Test store creation
store = mlflow.tracking._tracking_service.utils._get_store('modelregistry://localhost:8080')
print('Store type:', type(store).__name__)

# Test CA certificate configuration
import logging
logging.getLogger("model_registry_mlflow.api_client").setLevel(logging.INFO)

# This will show CA certificate detection logs
store = mlflow.tracking._tracking_service.utils._get_store('modelregistry+https://your-registry:8080')
```

**Testing CA Certificate with curl:**

```bash
# Test with custom CA certificate
curl --cacert /path/to/ca.crt https://your-registry:8080/api/model_registry/v1alpha3

# Test with system CA
curl https://your-registry:8080/api/model_registry/v1alpha3
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `uv run pytest`
6. Build and test the package: `uv build && uv pip install dist/*.whl --force-reinstall`
7. Submit a pull request

## License

This project uses the same license as the parent Model Registry project.

## Related Projects

- [Kubeflow Model Registry](https://github.com/kubeflow/model-registry)
- [MLflow](https://mlflow.org/)
- [Kubeflow](https://kubeflow.org/)