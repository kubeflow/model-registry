#!/usr/bin/env bash

set -e

MR_NAMESPACE="${MR_NAMESPACE:-kubeflow}"

if [[ -n "$LOCAL" ]]; then
    CLUSTER_NAME="${CLUSTER_NAME:-kind}"
    IMG="${IMG:-docker.io/kubeflow/model-registry:main}"

    echo 'Creating local Kind cluster and loading image'
    kind create cluster -n "$CLUSTER_NAME"
    kind load docker-image -n "$CLUSTER_NAME" "$IMG"
    echo 'Image loaded into kind cluster - use this command to port forward the mr service:'
    echo "kubectl port-forward -n $MR_NAMESPACE service/model-registry-service 8080:8080 &"
fi

echo 'Deploying model registry to Kind cluster'

kubectl create namespace "$MR_NAMESPACE"
kubectl apply -k manifests/kustomize/overlays/db
kubectl set image -n kubeflow deployment/model-registry-deployment rest-container="$IMG"
kubectl wait --for=condition=available -n "$MR_NAMESPACE" deployment/model-registry-db --timeout=5m
kubectl wait --for=condition=available -n "$MR_NAMESPACE" deployment/model-registry-deployment --timeout=5m
