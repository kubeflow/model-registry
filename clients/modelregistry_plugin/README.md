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

## Installation

```bash
pip install modelregistry_plugin
```

## Configuration

### Environment Variables

Set the following environment variables to configure the plugin:

```bash
export MLFLOW_TRACKING_URI="modelregistry://your-model-registry-host:port"
export MODEL_REGISTRY_HOST="your-model-registry-host"
export MODEL_REGISTRY_PORT="8080"
export MODEL_REGISTRY_SECURE="true"  # Use HTTPS
```

### Authentication

The plugin supports various authentication methods:

```bash
# Token-based auth
export MODEL_REGISTRY_TOKEN="your-token"

# Kubernetes service account token from`/var/run/secrets/kubernetes.io/serviceaccount/token`
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
    host="model-registry.kubeflow.svc.cluster.local",
    port=8080,
    secure=True,
    token="your-auth-token"
)

mlflow.set_tracking_uri(store.get_tracking_uri())
```

## Plugin Architecture

The plugin implements MLflow's tracking store interface by:

1. **Store Implementation**: `ModelRegistryStore` class implements `AbstractStore`
2. **API Translation**: Converts MLflow API calls to Model Registry API requests
3. **Authentication**: Handles Kubeflow authentication mechanisms
4. **Artifact Management**: Manages artifact storage and retrieval

## Supported Operations

- ✅ Create/list/get experiments
- ✅ Create/update/delete runs
- ✅ Log parameters, metrics, and tags
- ✅ Log and retrieve artifacts
- ✅ Search experiments and runs
- ✅ Model registration and versioning

## Development

### Setup Development Environment

```bash
# Clone repository
git clone <repository-url>
cd clients/modelregistry_plugin

# Install dependencies
uv sync

# Install in development mode
pip install -e .
```

### Running Tests

```bash
# Run all tests
uv run pytest

# Run with coverage
uv run pytest --cov=modelregistry_plugin
```

### Testing with Model Registry

```bash
# Start Model Registry locally
make start/mysql
make run/proxy

# Set test environment
export MLFLOW_TRACKING_URI="modelregistry://localhost:8080"

# Run integration tests
uv run pytest tests/integration/
```

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure Model Registry is running and accessible
2. **Authentication Failed**: Verify tokens and credentials are correct
3. **SSL/TLS Errors**: Check certificate configuration for secure connections

### Debug Logging

Enable debug logging to troubleshoot issues:

```python
import logging
logging.getLogger("modelregistry_plugin").setLevel(logging.DEBUG)
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project uses the same license as the parent Model Registry project.

## Related Projects

- [Kubeflow Model Registry](https://github.com/kubeflow/model-registry)
- [MLflow](https://mlflow.org/)
- [Kubeflow](https://kubeflow.org/)