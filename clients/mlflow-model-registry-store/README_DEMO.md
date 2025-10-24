# MLflow Model Registry Demo Script

This comprehensive demo script showcases all MLflow model registry functionality and can be run against either the native MLflow backend or the Model Registry backend.

## Demo Script Features

The `demo_mlflow_registry.py` script demonstrates:

### Core Model Registry Operations
- **Model Registration**: During logging and after logging
- **Model Management**: Creating, updating, versioning
- **Tags and Metadata**: Model and version tagging, descriptions
- **Aliases**: Setting and retrieving model aliases (champion, production, staging)
- **Model Loading**: By version number, latest version, and aliases
- **Search and Discovery**: Finding models and versions
- **Model Deletion**: Soft delete operations
- **Error Handling**: Proper exception handling for missing resources

### Supported MLflow Functions
All major MLflow model registry functions are demonstrated:

```python
# Model Registration
mlflow.sklearn.log_model(registered_model_name="model-name")
mlflow.register_model("runs:/run_id/model", "model-name")

# Model Management
client.create_registered_model("model-name")
client.create_model_version(name, source, run_id)
client.update_registered_model(name, description)
client.update_model_version(name, version, description)

# Tags and Metadata
client.set_registered_model_tag(name, key, value)
client.set_model_version_tag(name, version, key, value)
client.delete_registered_model_tag(name, key)
client.delete_model_version_tag(name, version, key)

# Aliases
client.set_registered_model_alias(name, alias, version)
client.get_model_version_by_alias(name, alias)
client.delete_registered_model_alias(name, alias)

# Model Loading
mlflow.sklearn.load_model("models:/model-name/1")
mlflow.sklearn.load_model("models:/model-name/latest")
mlflow.sklearn.load_model("models:/model-name@champion")

# Search Operations
client.search_registered_models()
client.search_model_versions(filter_string)
client.get_latest_versions(name)

# Deletion Operations
client.delete_model_version(name, version)
client.delete_registered_model(name)
```

## Usage

### Run Against Native MLflow Backend (Default)
```bash
python demo_mlflow_registry.py
```

This uses a temporary SQLite database for the MLflow backend.

### Run Against Model Registry Backend
```bash
# Set the Model Registry URI
export MLFLOW_REGISTRY_URI="modelregistry://localhost:8080?author=demo&is-secure=false"
python demo_mlflow_registry.py
```

Or inline:
```bash
MLFLOW_REGISTRY_URI="modelregistry://localhost:8080?author=demo&is-secure=false" python demo_mlflow_registry.py
```

### Example Model Registry URIs
```bash
# Local development
MLFLOW_REGISTRY_URI="modelregistry://localhost:8080?author=demo&is-secure=false"

# Kubernetes cluster
MLFLOW_REGISTRY_URI="modelregistry://model-registry.example.com:443?author=myuser&user-token=/path/to/token"

# With custom CA
MLFLOW_REGISTRY_URI="modelregistry://mr-server:8443?author=admin&custom-ca=/path/to/ca.crt"
```

## Expected Output

The demo will show:

1. **Setup Information**: Which backend is being used
2. **Model Registration**: Creating models during and after logging
3. **Management Operations**: Multiple versions with different configurations
4. **Tags and Metadata**: Comprehensive tagging of models and versions
5. **Aliases**: Setting up champion, production, and staging aliases
6. **Model Loading**: Loading models by version and alias
7. **Search Operations**: Finding and listing models
8. **Deletion**: Soft delete operations (archive-based for Model Registry)
9. **Error Handling**: Proper handling of missing resources
10. **Success Summary**: Confirmation all operations completed

## Compatibility

The script is designed to work identically with both backends:

- **Native MLflow**: Full functionality using SQLite backend
- **Model Registry**: Full functionality using archive-based operations

### MLflow Architecture: Tracking vs Model Registry

MLflow separates **tracking** (experiments, runs, metrics) from **model registry** (registered models, versions, aliases):

- **Tracking Backend**: Always uses SQLite in this demo (creates temporary database)
- **Model Registry Backend**: Uses native MLflow (default) or Model Registry server (when `MLFLOW_REGISTRY_URI` is set)

**Note**: When running `demo-modelregistry.sh`, you'll see SQLite database creation messages. This is expected - MLflow still needs a tracking backend for experiment/run data even when using a separate Model Registry backend for model storage.

### Key Differences
- **Delete Operations**: Native MLflow hard deletes, Model Registry archives
- **Search**: Model Registry implements client-side filtering for archived models
- **Stages**: Native MLflow has built-in stages, Model Registry uses custom properties

## Requirements

```bash
# Install the MLflow Model Registry Store plugin
uv pip install -e .

# Required dependencies
pip install mlflow scikit-learn numpy
```

## Testing

Run the demo script to verify:
1. All MLflow model registry functions work correctly
2. The plugin provides full compatibility between backends
3. Error handling works as expected
4. Performance is acceptable for typical workflows

The demo serves as both a comprehensive test and a reference for MLflow model registry usage patterns.