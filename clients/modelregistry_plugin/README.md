# Kubeflow Model Registry MLflow Plugin

A MLflow tracking plugin that integrates with Kubeflow Model Registry, enabling seamless experiment tracking and model management through MLflow's familiar API while storing metadata in the Kubeflow Model Registry backend.

## Overview

This plugin allows you to use MLflow's tracking API to log experiments, runs, metrics, parameters, and artifacts while having all metadata stored in Kubeflow Model Registry instead of MLflow's default backend storage.

## Features

- **Seamless Integration**: Use standard MLflow tracking APIs with Model Registry backend
- **Experiment Management**: Create and manage experiments in Model Registry
- **Run Tracking**: Log runs, metrics, parameters, and tags
- **Artifact Storage**: Store and retrieve model artifacts
- **Authentication**: Built-in support for Kubeflow authentication
- **Full MLflow Compatibility**: Implements all required AbstractStore methods

## Installation

### Using uv (Recommended)

```bash
# Build from source
cd clients/modelregistry_plugin
uv build
uv pip install dist/modelregistry_plugin-0.1.0-py3-none-any.whl
```

### Using pip

```bash
pip install modelregistry_plugin
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
from modelregistry_plugin import ModelRegistryStore

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
modelregistry = "modelregistry_plugin.tracking_store:ModelRegistryStore"
```

This allows MLflow to automatically discover and use the plugin when the `modelregistry://` URI scheme is specified.

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure Model Registry is running and accessible
2. **Authentication Failed**: Verify tokens and credentials are correct
3. **SSL/TLS Errors**: Check certificate configuration for secure connections
4. **Entry Point Not Found**: Ensure the package is properly installed and the entry point is registered

### Debug Logging

Enable debug logging to troubleshoot issues:

```python
import logging
logging.getLogger("modelregistry_plugin").setLevel(logging.DEBUG)
```

### Verification Commands

```python
# Check if the plugin is registered
import mlflow
print('Available tracking stores:', list(mlflow.tracking._tracking_service.utils._tracking_store_registry._registry.keys()))

# Test store creation
store = mlflow.tracking._tracking_service.utils._get_store('modelregistry://localhost:8080')
print('Store type:', type(store).__name__)
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