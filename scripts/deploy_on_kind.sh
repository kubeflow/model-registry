#!/usr/bin/env bash

set -e

MR_NAMESPACE="${MR_NAMESPACE:-kubeflow}"

# modularity to allow re-use this script against a remote k8s cluster
if [[ -n "$LOCAL" ]]; then
    CLUSTER_NAME="${CLUSTER_NAME:-kind}"
    IMG="${IMG:-kubeflow/model-registry:latest}"

    echo 'Creating local Kind cluster and loading image'

    if [[ $(kind get clusters || false) =~ $CLUSTER_NAME ]]; then
        echo 'Cluster already exists, skipping creation'

        kubectl config use-context "kind-$CLUSTER_NAME"
    else
        kind create cluster -n "$CLUSTER_NAME"
    fi

    kind load docker-image -n "$CLUSTER_NAME" "$IMG"

    echo 'Image loaded into kind cluster - use this command to port forward the mr service:'
    echo "kubectl port-forward -n $MR_NAMESPACE service/model-registry-service 8080:8080 &"
fi

echo 'Deploying model registry to Kind cluster'
if [[ $(kubectl get namespaces || false) =~ $MR_NAMESPACE ]]; then
    echo 'Namespace already exists, skipping creation'
else
    kubectl create namespace "$MR_NAMESPACE"
fi

kubectl apply -k manifests/kustomize/overlays/db
kubectl set image -n kubeflow deployment/model-registry-deployment rest-container="$IMG"
kubectl wait --for=condition=available -n "$MR_NAMESPACE" deployment/model-registry-db --timeout=5m
kubectl wait --for=condition=available -n "$MR_NAMESPACE" deployment/model-registry-deployment --timeout=5m
