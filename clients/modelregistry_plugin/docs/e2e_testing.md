# End-to-End Testing Guide

This guide explains how to set up and run end-to-end tests for the Model Registry MLflow plugin.

## Overview

The e2e tests verify that the `ModelRegistryStore` can successfully connect to a real Model Registry server and perform all supported operations. These tests are essential for ensuring the plugin works correctly in real-world scenarios.

## Prerequisites

### 1. Model Registry Server

You need access to a running Model Registry server. This could be:

- **Local Development**: Model Registry running on your local machine
- **Kubernetes Cluster**: Model Registry deployed in a K8s cluster
- **Cloud Instance**: Model Registry hosted in the cloud

### 2. Authentication

You need a valid authentication token for the Model Registry server. The token should have sufficient permissions to:

- Create and manage experiments
- Create and manage runs
- Log metrics, parameters, and artifacts
- Search and retrieve data

### 3. Network Access

Ensure you can reach the Model Registry server from your test environment:

```bash
# Test basic connectivity
curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://YOUR_HOST:PORT/health"
```

## Configuration

### Environment Variables

Set the following environment variables:

```bash
# Required
export MODEL_REGISTRY_HOST="your-model-registry-server.com"
export MODEL_REGISTRY_TOKEN="your-auth-token"

# Optional (with defaults)
export MODEL_REGISTRY_PORT="8080"  # defaults to 8080
export MODEL_REGISTRY_SECURE="false"  # defaults to false (HTTP)
```

### Configuration File

Alternatively, create a configuration file:

```bash
# Copy the example
cp tests/e2e_config.env.example tests/e2e_config.env

# Edit with your values
nano tests/e2e_config.env

# Source the configuration
source tests/e2e_config.env
```

Example configuration file:

```bash
# Model Registry server details
MODEL_REGISTRY_HOST=model-registry.kubeflow.svc.cluster.local
MODEL_REGISTRY_PORT=8080
MODEL_REGISTRY_SECURE=false
MODEL_REGISTRY_TOKEN=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...

# Optional debugging
LOG_LEVEL=DEBUG
```

## Running Tests

### Quick Setup Test

Before running the full e2e test suite, verify your setup:

```bash
# Test environment and connectivity
python scripts/test_e2e_setup.py
```

This script will:
- Check environment variables
- Test ModelRegistryStore instantiation
- Verify connectivity to the server
- Test MLflow integration

### Full E2E Test Suite

**Option 1: Using the provided script (Recommended)**

```bash
# Make executable (first time only)
chmod +x scripts/run_e2e_tests.sh

# Run all e2e tests
./scripts/run_e2e_tests.sh
```

**Option 2: Manual execution**

```bash
# Set environment variables
export MODEL_REGISTRY_HOST="your-server.com"
export MODEL_REGISTRY_TOKEN="your-token"

# Run e2e tests
uv run pytest tests/test_e2e.py -v -s
```

**Option 3: With configuration file**

```bash
# Source configuration
source tests/e2e_config.env

# Run tests
uv run pytest tests/test_e2e.py -v -s
```

## Test Coverage

The e2e tests cover the following functionality:

### Connection Tests
- ✅ Basic connectivity to Model Registry server
- ✅ Authentication verification
- ✅ Server health check

### Experiment Management
- ✅ Create experiments
- ✅ Retrieve experiments by ID and name
- ✅ Search experiments with different view types
- ✅ Experiment lifecycle (delete/restore)
- ✅ Experiment tags

### Run Management
- ✅ Create runs
- ✅ Retrieve runs by ID
- ✅ Update run information and status
- ✅ Run lifecycle (delete/restore)
- ✅ Run tags

### Data Logging
- ✅ Log individual metrics and parameters
- ✅ Batch logging of metrics, parameters, and tags
- ✅ Metric history retrieval
- ✅ Tag management (set/delete)

### Search and Filtering
- ✅ Search experiments with pagination
- ✅ Search runs with different view types
- ✅ Filter by experiment IDs

### MLflow Integration
- ✅ MLflow tracking URI configuration
- ✅ MLflow API compatibility
- ✅ Entry point registration verification

## Test Output

### Successful Run

```
🚀 Model Registry E2E Setup Test
==================================================
🔍 Checking environment variables...
  ✅ MODEL_REGISTRY_HOST: your-server.com
  ✅ MODEL_REGISTRY_TOKEN: your-token...
  ℹ️  MODEL_REGISTRY_PORT: 8080
  ℹ️  MODEL_REGISTRY_SECURE: false

🔍 Testing ModelRegistryStore instantiation...
  📡 Store URI: modelregistry://your-server.com:8080
  ✅ ModelRegistryStore instantiated successfully: ModelRegistryStore

🔍 Testing connectivity to Model Registry server...
  ✅ Successfully connected to Model Registry server
  📊 Found 5 experiments

🔍 Testing MLflow integration...
  📋 Available tracking stores: ['file', 'sqlite', 'modelregistry', ...]
  ✅ modelregistry tracking store is registered with MLflow
  ✅ Successfully set MLflow tracking URI: modelregistry://your-server.com:8080

==================================================
✅ All tests passed! E2E setup is ready.
```

### Test Results

```
[INFO] Running E2E tests...
test_store_connection PASSED
test_create_and_get_experiment PASSED
test_create_and_get_run PASSED
test_log_metrics_and_params PASSED
test_log_batch PASSED
test_search_experiments PASSED
test_search_runs PASSED
test_experiment_tags PASSED
test_run_tags PASSED
test_experiment_lifecycle PASSED
test_run_lifecycle PASSED
test_mlflow_integration PASSED

[SUCCESS] All E2E tests passed!
```

## Troubleshooting

### Common Issues

#### 1. Connection Refused

**Symptoms**: `ConnectionError` or `ConnectionRefusedError`

**Solutions**:
- Verify the Model Registry server is running
- Check the host and port are correct
- Ensure network connectivity (firewall, VPN, etc.)
- Test with curl: `curl http://HOST:PORT/health`

#### 2. Authentication Failed

**Symptoms**: `401 Unauthorized` or `403 Forbidden`

**Solutions**:
- Verify the token is valid and not expired
- Check token permissions
- Ensure token format is correct
- Test with curl: `curl -H "Authorization: Bearer TOKEN" http://HOST:PORT/health`

#### 3. SSL/TLS Errors

**Symptoms**: `SSLError` or certificate verification failures

**Solutions**:
- Set `MODEL_REGISTRY_SECURE=true` for HTTPS
- Check certificate configuration
- Verify server certificate is valid
- For self-signed certificates, you may need to configure certificate handling

#### 4. Entry Point Not Found

**Symptoms**: `modelregistry` not in available tracking stores

**Solutions**:
- Ensure the package is properly installed
- Rebuild and reinstall: `uv build && uv pip install dist/*.whl --force-reinstall`
- Check entry point registration in `pyproject.toml`

### Debug Mode

Enable debug logging for more detailed information:

```bash
# Set debug level
export LOG_LEVEL=DEBUG

# Run with debug output
uv run pytest tests/test_e2e.py -v -s --log-cli-level=DEBUG
```

### Manual Testing

If automated tests fail, you can test manually:

```python
import os
from modelregistry_plugin.store import ModelRegistryStore

# Set environment
os.environ["MODEL_REGISTRY_TOKEN"] = "your-token"

# Create store
store = ModelRegistryStore(store_uri="modelregistry://your-host:8080")

# Test basic operations
experiments = store.list_experiments()
print(f"Found {len(experiments)} experiments")
```

## Continuous Integration

For CI/CD pipelines, you can run e2e tests with:

```yaml
# Example GitHub Actions workflow
- name: Run E2E Tests
  env:
    MODEL_REGISTRY_HOST: ${{ secrets.MODEL_REGISTRY_HOST }}
    MODEL_REGISTRY_TOKEN: ${{ secrets.MODEL_REGISTRY_TOKEN }}
    MODEL_REGISTRY_PORT: ${{ secrets.MODEL_REGISTRY_PORT }}
  run: |
    ./scripts/run_e2e_tests.sh
```

## Best Practices

1. **Use Dedicated Test Environment**: Run e2e tests against a dedicated test instance
2. **Clean Up After Tests**: Tests should clean up their own data
3. **Use Unique Names**: Use UUIDs or timestamps to avoid conflicts
4. **Monitor Test Duration**: E2e tests can be slow; monitor and optimize
5. **Handle Flaky Tests**: Network issues can cause intermittent failures
6. **Log Test Artifacts**: Save logs and artifacts for debugging

## Security Considerations

- **Token Security**: Never commit tokens to version control
- **Network Security**: Use secure connections (HTTPS) in production
- **Access Control**: Use tokens with minimal required permissions
- **Audit Logging**: Monitor test activities in production environments 