#!/usr/bin/env python3
"""
Simple script to test E2E setup and verify connectivity to Model Registry server.

This script can be used to verify that:
1. Environment variables are properly set
2. The ModelRegistryStore can be instantiated
3. Basic connectivity to the Model Registry server works

Usage:
    python scripts/test_e2e_setup.py
"""

import os
import sys
from pathlib import Path

# Add the project root to the Python path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root))

def check_environment():
    """Check if required environment variables are set."""
    print("🔍 Checking environment variables...")
    
    required_vars = ["MODEL_REGISTRY_HOST", "MODEL_REGISTRY_TOKEN"]
    optional_vars = ["MODEL_REGISTRY_PORT", "MODEL_REGISTRY_SECURE"]
    
    missing_vars = []
    for var in required_vars:
        if not os.getenv(var):
            missing_vars.append(var)
        else:
            print(f"  ✅ {var}: {os.getenv(var)[:10]}..." if var == "MODEL_REGISTRY_TOKEN" else f"  ✅ {var}: {os.getenv(var)}")
    
    for var in optional_vars:
        value = os.getenv(var, "not set (using default)")
        print(f"  ℹ️  {var}: {value}")
    
    if missing_vars:
        print(f"  ❌ Missing required environment variables: {', '.join(missing_vars)}")
        return False
    
    return True

def test_store_instantiation():
    """Test that ModelRegistryStore can be instantiated."""
    print("\n🔍 Testing ModelRegistryStore instantiation...")
    
    try:
        from modelregistry_plugin.store import ModelRegistryStore
        
        host = os.getenv("MODEL_REGISTRY_HOST")
        port = os.getenv("MODEL_REGISTRY_PORT", "8080")
        secure = os.getenv("MODEL_REGISTRY_SECURE", "false").lower() == "true"
        
        store_uri = f"modelregistry://{host}:{port}"
        print(f"  📡 Store URI: {store_uri}")
        
        store = ModelRegistryStore(store_uri=store_uri)
        print(f"  ✅ ModelRegistryStore instantiated successfully: {type(store).__name__}")
        
        return store
        
    except Exception as e:
        print(f"  ❌ Failed to instantiate ModelRegistryStore: {e}")
        return None

def test_connectivity(store):
    """Test basic connectivity to the Model Registry server."""
    print("\n🔍 Testing connectivity to Model Registry server...")
    
    try:
        # Try to list experiments to verify connection
        experiments = store.list_experiments()
        print(f"  ✅ Successfully connected to Model Registry server")
        print(f"  📊 Found {len(experiments)} experiments")
        
        # Show some experiment details if available
        if experiments:
            print("  📋 Sample experiments:")
            for exp in experiments[:3]:  # Show first 3 experiments
                print(f"    - {exp.name} (ID: {exp.experiment_id})")
            if len(experiments) > 3:
                print(f"    ... and {len(experiments) - 3} more")
        
        return True
        
    except Exception as e:
        print(f"  ❌ Failed to connect to Model Registry server: {e}")
        return False

def test_mlflow_integration():
    """Test MLflow integration."""
    print("\n🔍 Testing MLflow integration...")
    
    try:
        import mlflow
        
        # Check if modelregistry is available as a tracking store
        tracking_stores = list(mlflow.tracking._tracking_service.utils._tracking_store_registry._registry.keys())
        print(f"  📋 Available tracking stores: {tracking_stores}")
        
        if "modelregistry" in tracking_stores:
            print("  ✅ modelregistry tracking store is registered with MLflow")
            
            # Test setting tracking URI
            host = os.getenv("MODEL_REGISTRY_HOST")
            port = os.getenv("MODEL_REGISTRY_PORT", "8080")
            tracking_uri = f"modelregistry://{host}:{port}"
            
            mlflow.set_tracking_uri(tracking_uri)
            print(f"  ✅ Successfully set MLflow tracking URI: {tracking_uri}")
            
            return True
        else:
            print("  ❌ modelregistry tracking store is not registered with MLflow")
            return False
            
    except Exception as e:
        print(f"  ❌ Failed to test MLflow integration: {e}")
        return False

def main():
    """Main function to run all tests."""
    print("🚀 Model Registry E2E Setup Test")
    print("=" * 50)
    
    # Check environment
    if not check_environment():
        print("\n❌ Environment check failed. Please set the required environment variables.")
        print("\nExample:")
        print("  export MODEL_REGISTRY_HOST='your-server.com'")
        print("  export MODEL_REGISTRY_TOKEN='your-token'")
        print("  export MODEL_REGISTRY_PORT='8080'  # optional")
        print("  export MODEL_REGISTRY_SECURE='false'  # optional")
        sys.exit(1)
    
    # Test store instantiation
    store = test_store_instantiation()
    if not store:
        print("\n❌ Store instantiation failed.")
        sys.exit(1)
    
    # Test connectivity
    if not test_connectivity(store):
        print("\n❌ Connectivity test failed.")
        sys.exit(1)
    
    # Test MLflow integration
    if not test_mlflow_integration():
        print("\n❌ MLflow integration test failed.")
        sys.exit(1)
    
    print("\n" + "=" * 50)
    print("✅ All tests passed! E2E setup is ready.")
    print("\nYou can now run the full E2E test suite:")
    print("  ./scripts/run_e2e_tests.sh")
    print("  # or")
    print("  uv run pytest tests/test_e2e.py -v -s")

if __name__ == "__main__":
    main() 