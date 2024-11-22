#!/usr/bin/env bash

set -e

DIR="$(dirname "$0")"
MR_NAMESPACE="${MR_NAMESPACE:-kubeflow}"
IMG="${IMG:-kubeflow/model-registry:latest}"

source ./${DIR}/utils.sh

# modularity to allow re-use this script against a remote k8s cluster
if [[ -n "$LOCAL" ]]; then
    CLUSTER_NAME="${CLUSTER_NAME:-kind}"

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
kubectl patch deployment -n "$MR_NAMESPACE" model-registry-deployment \
--patch '{"spec": {"template": {"spec": {"containers": [{"name": "rest-container", "image": "'$IMG'", "imagePullPolicy": "IfNotPresent"}]}}}}'

kubectl wait --for=condition=available -n "$MR_NAMESPACE" deployment/model-registry-db --timeout=5m

kubectl delete pod -n "$MR_NAMESPACE" --selector='component=model-registry-server'

repeat_cmd_until "kubectl get pod -n "$MR_NAMESPACE" --selector='component=model-registry-server' \
-o jsonpath=\"{.items[*].spec.containers[?(@.name=='rest-container')].image}\"" "= $IMG" 300 "kubectl describe pod -n $MR_NAMESPACE --selector='component=model-registry-server'"

kubectl wait --for=condition=available -n "$MR_NAMESPACE" deployment/model-registry-deployment --timeout=5m
