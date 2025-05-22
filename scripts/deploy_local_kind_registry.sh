#!/bin/bash

source "$(cd "$(dirname "$0")" && pwd)/common.sh"

SCRIPT_DIR="$(dirname "$(realpath_fallback "$BASH_SOURCE")")"
set -e

kubectl apply -f "${SCRIPT_DIR}/services/container_registry.yaml"

echo "Waiting for Deployment..."
kubectl wait --for=condition=available deployment/distribution-registry-test-deployment --timeout=5m
kubectl logs deployment/distribution-registry-test-deployment
echo "Deployment looks ready."

