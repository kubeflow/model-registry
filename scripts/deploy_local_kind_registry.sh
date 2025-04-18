#!/bin/bash

SCRIPT_DIR="$(dirname "$(realpath "$BASH_SOURCE")")"
set -e

kubectl apply -f "${SCRIPT_DIR}/services/container_registry.yaml"

echo "Waiting for Deployment..."
kubectl wait --for=condition=available deployment/distribution-registry-test-deployment --timeout=5m
kubectl logs deployment/distribution-registry-test-deployment
echo "Deployment looks ready."

