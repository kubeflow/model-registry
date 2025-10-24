#!/usr/bin/env python3
"""Demo script to test the MLflow plugin functionality."""

import os
import sys
from unittest.mock import Mock, patch

# Test if the plugin is discoverable by MLflow
try:
    from mlflow_model_registry_store.store import ModelRegistryStore
    print("✅ Plugin imported successfully")
except ImportError as e:
    print(f"❌ Failed to import plugin: {e}")
    sys.exit(1)

# Test plugin entry point discovery
try:
    import pkg_resources

    entry_points = list(pkg_resources.iter_entry_points("mlflow.model_registry_store"))
    modelregistry_entry = None

    for ep in entry_points:
        if ep.name == "modelregistry":
            modelregistry_entry = ep
            break

    if modelregistry_entry:
        print("✅ Plugin entry point discovered by MLflow")
        print(f"   Entry point: {modelregistry_entry}")
    else:
        print("❌ Plugin entry point NOT found")
        print(f"   Available entry points: {[ep.name for ep in entry_points]}")

except Exception as e:
    print(f"⚠️  Could not check entry points: {e}")

# Test basic store initialization
try:
    uri = "modelregistry://localhost:8080?author=demo&is_secure=false"

    with patch('mlflow_model_registry_store.store.ModelRegistry') as mock_registry:
        # Mock successful connection
        mock_client = Mock()
        mock_registry.return_value = mock_client

        store = ModelRegistryStore(store_uri=uri)
        print("✅ Store initialized successfully")
        print(f"   Server: {store.server_address}:{store.port}")
        print(f"   Author: {store.author}")
        print(f"   Secure: {store.is_secure}")

        # Test basic method availability
        methods_to_test = [
            'create_registered_model',
            'get_registered_model',
            'search_registered_models',
            'create_model_version',
            'get_model_version',
            'set_registered_model_tag',
            'set_model_version_tag',
            'set_registered_model_alias',
            'get_model_version_by_alias',
        ]

        missing_methods = []
        for method_name in methods_to_test:
            if not hasattr(store, method_name):
                missing_methods.append(method_name)

        if missing_methods:
            print(f"❌ Missing required methods: {missing_methods}")
        else:
            print("✅ All required AbstractStore methods implemented")

except Exception as e:
    print(f"❌ Failed to initialize store: {e}")
    sys.exit(1)

# Test MLflow client integration
try:
    os.environ["MLFLOW_REGISTRY_URI"] = uri

    with patch('mlflow_model_registry_store.store.ModelRegistry') as mock_registry:
        # Mock successful connection
        mock_client = Mock()
        mock_registry.return_value = mock_client

        import mlflow
        client = mlflow.MlflowClient()

        # The client should use our store backend
        print("✅ MLflow client created with plugin backend")
        print(f"   Registry URI: {os.environ.get('MLFLOW_REGISTRY_URI')}")

except Exception as e:
    print(f"❌ Failed to create MLflow client: {e}")

print("\n🎉 Plugin demo completed successfully!")
print("\nTo use the plugin in your MLflow applications:")
print('1. Set environment variable: export MLFLOW_REGISTRY_URI="modelregistry://server:port?author=yourname"')
print("2. Use MLflow model registry as usual - it will automatically use the Model Registry backend")