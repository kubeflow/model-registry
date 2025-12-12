#!/bin/bash

SCRIPT_DIR="$(dirname "$(realpath "$BASH_SOURCE")")"
OCI_REGISTRY_NAMESPACE="local-oci-registry-ns"
set -e
if [[ $(kubectl get namespaces || false) =~ $OCI_REGISTRY_NAMESPACE ]]; then
    echo 'Namespace already exists, skipping creation'
else
    kubectl create namespace "$OCI_REGISTRY_NAMESPACE"
fi

kubectl apply -f "${SCRIPT_DIR}/services/container_registry.yaml" -n $OCI_REGISTRY_NAMESPACE

echo "Waiting for Deployment..."
kubectl wait --for=condition=available deployment/distribution-registry-test-deployment -n $OCI_REGISTRY_NAMESPACE --timeout=5m
kubectl logs deployment/distribution-registry-test-deployment -n $OCI_REGISTRY_NAMESPACE
echo "Deployment looks ready."

