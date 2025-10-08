#!/bin/bash
# Run MLflow Model Registry Demo with Model Registry Backend
exec env MLFLOW_REGISTRY_URI="modelregistry://localhost:8080?author=demo&is-secure=false" uv run python demo_mlflow_registry.py