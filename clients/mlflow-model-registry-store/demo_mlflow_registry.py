#!/usr/bin/env python3
"""
MLflow Model Registry End-to-End Demo Script

This script demonstrates all MLflow model registry functionality and can be run against:
1. Native MLflow backend (SQLite) - default
2. Model Registry backend - set MLFLOW_REGISTRY_URI=modelregistry://...

Usage:
    # Run against native MLflow backend
    python demo_mlflow_registry.py

    # Run against Model Registry backend
    MLFLOW_REGISTRY_URI="modelregistry://localhost:8080?author=demo&is-secure=false" python demo_mlflow_registry.py
"""

import logging
import os
import tempfile
import time
import uuid
import warnings
from pathlib import Path

# Configure logging to suppress specific MLflow warnings
class MLflowWarningFilter(logging.Filter):
    """Custom filter to suppress specific MLflow warnings."""

    def filter(self, record):
        # Filter out pip version warnings
        if "Failed to resolve installed pip version" in record.getMessage():
            return False
        if "pip`` will be added to conda.yaml environment spec without a version" in record.getMessage():
            return False
        # Filter out tracking registry warnings about missing artifacts
        if "has no artifacts at artifact path" in record.getMessage():
            return False
        return True

# Apply the filter to MLflow loggers
mlflow_loggers = [
    'mlflow.utils.environment',
    'mlflow.tracking._model_registry.fluent',
    'mlflow.models.model',
    'mlflow'
]

for logger_name in mlflow_loggers:
    logger = logging.getLogger(logger_name)
    logger.addFilter(MLflowWarningFilter())
    # Set level to WARNING to allow important warnings but filter out noise
    logger.setLevel(logging.WARNING)

# Suppress Python warnings from MLflow
warnings.filterwarnings("ignore", message=".*Failed to resolve installed pip version.*")
warnings.filterwarnings("ignore", message=".*get_latest_versions.*is deprecated.*", category=FutureWarning)
warnings.filterwarnings("ignore", category=UserWarning, module="mlflow.utils.environment")
warnings.filterwarnings("ignore", category=UserWarning, module="mlflow.tracking._model_registry.fluent")

import mlflow
import mlflow.sklearn
import numpy as np
from mlflow import MlflowClient
from sklearn.datasets import make_classification, make_regression
from sklearn.ensemble import RandomForestClassifier, RandomForestRegressor
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import accuracy_score, mean_squared_error
from sklearn.model_selection import train_test_split

def setup_tracking():
    """Set up MLflow tracking with temporary directory."""
    temp_dir = tempfile.mkdtemp()
    tracking_uri = f"sqlite:///{temp_dir}/mlflow.db"
    mlflow.set_tracking_uri(tracking_uri)

    # MLflow automatically reads MLFLOW_REGISTRY_URI environment variable
    registry_uri = mlflow.get_registry_uri()
    print(f"Registry backend: {registry_uri}")

    return temp_dir

def create_sample_data():
    """Create sample datasets for demo."""
    print("\n=== Creating Sample Data ===")

    # Regression dataset
    X_reg, y_reg = make_regression(
        n_samples=1000, n_features=10, n_informative=5,
        noise=0.1, random_state=42
    )
    X_reg_train, X_reg_test, y_reg_train, y_reg_test = train_test_split(
        X_reg, y_reg, test_size=0.2, random_state=42
    )

    # Classification dataset
    X_clf, y_clf = make_classification(
        n_samples=1000, n_features=20, n_informative=10, n_redundant=10,
        n_clusters_per_class=1, random_state=42
    )
    X_clf_train, X_clf_test, y_clf_train, y_clf_test = train_test_split(
        X_clf, y_clf, test_size=0.2, random_state=42
    )

    print(f"Created regression dataset: {X_reg.shape}")
    print(f"Created classification dataset: {X_clf.shape}")

    return (X_reg_train, X_reg_test, y_reg_train, y_reg_test,
            X_clf_train, X_clf_test, y_clf_train, y_clf_test)

def demo_model_registration_during_logging():
    """Demo: Model registration during logging."""
    print("\n=== Demo: Model Registration During Logging ===")

    # Get sample data
    data = create_sample_data()
    X_reg_train, X_reg_test, y_reg_train, y_reg_test = data[:4]

    model_name = f"regression-model-{uuid.uuid4().hex[:8]}"

    with mlflow.start_run(run_name="regression-training") as run:
        # Train model
        params = {"n_estimators": 100, "max_depth": 5, "random_state": 42}
        model = RandomForestRegressor(**params)
        model.fit(X_reg_train, y_reg_train)

        # Log parameters and metrics
        mlflow.log_params(params)
        y_pred = model.predict(X_reg_test)
        mse = mean_squared_error(y_reg_test, y_pred)
        mlflow.log_metric("mse", mse)

        # Log and register model in one step
        model_info = mlflow.sklearn.log_model(
            sk_model=model,
            name="model",
            input_example=X_reg_train[:5],
            registered_model_name=model_name
        )

        print(f"Registered model: {model_name}")
        print(f"Model version: {model_info.registered_model_version}")
        print(f"Run ID: {run.info.run_id}")

    return model_name, run.info.run_id

def demo_model_registration_after_logging():
    """Demo: Model registration after logging."""
    print("\n=== Demo: Model Registration After Logging ===")

    # Get sample data
    data = create_sample_data()
    X_clf_train, X_clf_test, y_clf_train, y_clf_test = data[4:]

    model_name = f"classification-model-{uuid.uuid4().hex[:8]}"

    with mlflow.start_run(run_name="classification-training") as run:
        # Train model
        params = {"C": 1.0, "random_state": 42}
        model = LogisticRegression(**params)
        model.fit(X_clf_train, y_clf_train)

        # Log parameters and metrics
        mlflow.log_params(params)
        y_pred = model.predict(X_clf_test)
        accuracy = accuracy_score(y_clf_test, y_pred)
        mlflow.log_metric("accuracy", accuracy)

        # Log model without registration
        mlflow.sklearn.log_model(
            sk_model=model,
            name="model",
            input_example=X_clf_train[:5]
        )

        run_id = run.info.run_id

    # Register model after logging
    model_uri = f"runs:/{run_id}/model"
    result = mlflow.register_model(model_uri, model_name)

    print(f"Registered model: {model_name}")
    print(f"Model version: {result.version}")
    print(f"Source: {model_uri}")

    return model_name, run_id

def demo_model_management():
    """Demo: model management operations."""
    print("\n=== Demo: Model Management Operations ===")

    client = MlflowClient()

    # Create a dedicated model for management demo
    management_model = f"management-demo-{uuid.uuid4().hex[:8]}"

    # Create registered model explicitly
    registered_model = client.create_registered_model(
        management_model,
        description="Demo model for management operations"
    )
    print(f"Created registered model: {management_model}")

    # Create multiple versions
    data = create_sample_data()
    X_train, X_test, y_train, y_test = data[:4]

    version_info = []
    for i in range(3):
        with mlflow.start_run(run_name=f"version-{i+1}") as run:
            # Train different model configurations
            params = {
                "n_estimators": 50 + i*25,
                "max_depth": 3 + i,
                "random_state": 42
            }
            model = RandomForestRegressor(**params)
            model.fit(X_train, y_train)

            mlflow.log_params(params)
            y_pred = model.predict(X_test)
            mse = mean_squared_error(y_test, y_pred)
            mlflow.log_metric("mse", mse)

            mlflow.sklearn.log_model(
                sk_model=model,
                name="model",
                input_example=X_train[:5]
            )

            # Create model version
            model_version = client.create_model_version(
                name=management_model,
                source=f"runs:/{run.info.run_id}/model",
                run_id=run.info.run_id,
                description=f"Version {i+1} with {params['n_estimators']} estimators"
            )

            version_info.append((model_version.version, mse))
            print(f"Created version {model_version.version} with MSE: {mse:.4f}")

    return management_model, version_info

def demo_tags_and_metadata():
    """Demo: Tags and metadata management."""
    print("\n=== Demo: Tags and Metadata Management ===")

    client = MlflowClient()

    # Use the management model from previous demo
    model_name, version_info = demo_model_management()

    # Set model-level tags
    client.set_registered_model_tag(model_name, "task", "regression")
    client.set_registered_model_tag(model_name, "team", "data-science")
    client.set_registered_model_tag(model_name, "project", "mlflow-demo")
    print(f"Set model tags for {model_name}")

    # Set version-level tags
    for version, mse in version_info:
        client.set_model_version_tag(
            model_name, version, "performance_tier",
            "high" if mse < 100 else "medium" if mse < 200 else "low"
        )
        client.set_model_version_tag(
            model_name, version, "validated", "true"
        )
        print(f"Set version tags for {model_name} v{version}")

    # Update model description
    client.update_registered_model(
        model_name,
        description="Updated: Regression model with multiple versions and tagging"
    )

    # Update version descriptions
    for version, mse in version_info:
        client.update_model_version(
            model_name, version,
            description=f"Updated: Version {version} - MSE: {mse:.4f}, Performance validated"
        )

    print("Updated model and version descriptions")

    return model_name

def demo_aliases():
    """Demo: Model aliases management."""
    print("\n=== Demo: Model Aliases Management ===")

    client = MlflowClient()
    model_name, version_info = demo_model_management()

    # Find best performing version
    best_version = min(version_info, key=lambda x: x[1])[0]
    latest_version = max(version_info, key=lambda x: int(x[0]))[0]

    # Set aliases
    client.set_registered_model_alias(model_name, "champion", best_version)
    client.set_registered_model_alias(model_name, "production", latest_version)
    client.set_registered_model_alias(model_name, "staging", "2")  # Middle version for staging

    print(f"Set alias 'champion' -> version {best_version}")
    print(f"Set alias 'production' -> version {latest_version}")
    print(f"Set alias 'staging' -> version 2")

    # Test alias retrieval
    champion_version = client.get_model_version_by_alias(model_name, "champion")
    print(f"Champion version: {champion_version.version}")

    return model_name

def demo_model_loading():
    """Demo: Different model loading patterns."""
    print("\n=== Demo: Model Loading Patterns ===")

    # Get models from previous demos
    reg_model, _ = demo_model_registration_during_logging()
    clf_model, _ = demo_model_registration_after_logging()
    alias_model = demo_aliases()

    # Wait a moment for models to be fully registered
    time.sleep(2)

    # Load by version number
    try:
        model_v1 = mlflow.sklearn.load_model(f"models:/{reg_model}/1")
        print(f"✓ Loaded {reg_model} version 1")
    except Exception as e:
        print(f"✗ Failed to load {reg_model} version 1: {e}")

    # Load latest version
    try:
        model_latest = mlflow.sklearn.load_model(f"models:/{clf_model}/latest")
        print(f"✓ Loaded {clf_model} latest version")
    except Exception as e:
        print(f"✗ Failed to load {clf_model} latest: {e}")

    # Load by alias
    try:
        model_champion = mlflow.sklearn.load_model(f"models:/{alias_model}@champion")
        print(f"✓ Loaded {alias_model} champion alias")
    except Exception as e:
        print(f"✗ Failed to load {alias_model} champion: {e}")

    # Load by production alias
    try:
        model_production = mlflow.sklearn.load_model(f"models:/{alias_model}@production")
        print(f"✓ Loaded {alias_model} production alias")
    except Exception as e:
        print(f"✗ Failed to load {alias_model} production: {e}")

    return reg_model, clf_model, alias_model

def demo_search_and_discovery():
    """Demo: Search and discovery operations."""
    print("\n=== Demo: Search and Discovery ===")

    client = MlflowClient()

    # Search registered models
    models = client.search_registered_models()
    print(f"Total registered models: {len(models)}")

    for model in models[:5]:  # Show first 5
        print(f"  - {model.name}: {model.description or 'No description'}")

    # Search model versions
    if models:
        model_name = models[0].name
        try:
            versions = client.search_model_versions(f"name='{model_name}'")
            print(f"Versions for {model_name}: {len(versions)}")
            for version in versions:
                print(f"  - Version {version.version}: {version.current_stage or 'No stage'}")
        except Exception as e:
            print(f"Search model versions not fully supported: {e}")

    # Get latest versions
    if models:
        try:
            latest_versions = client.get_latest_versions(models[0].name)
            print(f"Latest versions for {models[0].name}: {len(latest_versions)}")
        except Exception as e:
            print(f"Get latest versions error: {e}")

def demo_model_deletion():
    """Demo: Model deletion operations."""
    print("\n=== Demo: Model Deletion ===")

    client = MlflowClient()

    # Create a temporary model for deletion demo
    temp_model = f"delete-demo-{uuid.uuid4().hex[:8]}"

    # Create model and version
    client.create_registered_model(temp_model, description="Temporary model for deletion demo")

    with mlflow.start_run() as run:
        # Simple model for deletion
        data = create_sample_data()
        X_train, _, y_train, _ = data[:4]

        model = RandomForestRegressor(n_estimators=10, random_state=42)
        model.fit(X_train, y_train)
        mlflow.sklearn.log_model(
            sk_model=model,
            name="model",
            input_example=X_train[:5]
        )

        version = client.create_model_version(
            temp_model,
            f"runs:/{run.info.run_id}/model",
            run.info.run_id
        )

    print(f"Created temporary model: {temp_model} version {version.version}")

    # Delete model version
    try:
        client.delete_model_version(temp_model, version.version)
        print(f"✓ Deleted version {version.version}")

        # Verify version cannot be accessed after deletion
        try:
            client.get_model_version(temp_model, version.version)
            print(f"✗ ERROR: Version {version.version} should not be accessible after deletion")
        except Exception:
            print(f"✓ Verified: Version {version.version} is not accessible after deletion")

        # Verify version cannot be loaded after deletion
        try:
            mlflow.sklearn.load_model(f"models:/{temp_model}/{version.version}")
            print(f"✗ ERROR: Version {version.version} should not be loadable after deletion")
        except Exception:
            print(f"✓ Verified: Version {version.version} cannot be loaded after deletion")

    except Exception as e:
        print(f"Version deletion: {e}")

    # Delete registered model
    try:
        client.delete_registered_model(temp_model)
        print(f"✓ Deleted model {temp_model}")

        # Verify model cannot be accessed after deletion
        try:
            client.get_registered_model(temp_model)
            print(f"✗ ERROR: Model {temp_model} should not be accessible after deletion")
        except Exception:
            print(f"✓ Verified: Model {temp_model} is not accessible after deletion")

        # Verify model cannot be loaded after deletion
        try:
            mlflow.sklearn.load_model(f"models:/{temp_model}/latest")
            print(f"✗ ERROR: Model {temp_model} should not be loadable after deletion")
        except Exception:
            print(f"✓ Verified: Model {temp_model} cannot be loaded after deletion")

    except Exception as e:
        print(f"Model deletion: {e}")

def demo_error_handling():
    """Demo: Error handling for non-existent resources."""
    print("\n=== Demo: Error Handling ===")

    client = MlflowClient()

    # Test non-existent model
    try:
        client.get_registered_model("non-existent-model")
        print("✗ Should have failed for non-existent model")
    except Exception as e:
        print(f"✓ Correctly handled non-existent model: {type(e).__name__}")

    # Test non-existent version
    try:
        client.get_model_version("non-existent-model", "999")
        print("✗ Should have failed for non-existent version")
    except Exception as e:
        print(f"✓ Correctly handled non-existent version: {type(e).__name__}")

    # Test non-existent alias
    try:
        client.get_model_version_by_alias("non-existent-model", "non-existent-alias")
        print("✗ Should have failed for non-existent alias")
    except Exception as e:
        print(f"✓ Correctly handled non-existent alias: {type(e).__name__}")

def run_demo():
    """Run the complete MLflow model registry demo."""
    print("=" * 80)
    print("MLflow Model Registry Demo")
    print("=" * 80)

    # Setup
    temp_dir = setup_tracking()

    try:
        # Core registration demos
        demo_model_registration_during_logging()
        demo_model_registration_after_logging()

        # Management operations
        demo_tags_and_metadata()
        demo_aliases()

        # Advanced operations
        demo_model_loading()
        demo_search_and_discovery()
        demo_model_deletion()
        demo_error_handling()

        print("\n" + "=" * 80)
        print("✅ Demo completed successfully!")
        print("All MLflow model registry functions demonstrated.")
        print("=" * 80)

    except Exception as e:
        print(f"\n❌ Demo failed with error: {e}")
        import traceback
        traceback.print_exc()

    finally:
        # Cleanup
        import shutil
        try:
            shutil.rmtree(temp_dir, ignore_errors=True)
            print(f"Cleaned up temporary directory: {temp_dir}")
        except:
            pass

if __name__ == "__main__":
    run_demo()
