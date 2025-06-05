#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(realpath "$BASH_SOURCE")")"

echo "Check if Data Science Cluster exists"
DSC_NAME="default-dsc"

if kubectl get datasciencecluster "$DSC_NAME" &> /dev/null; then
  echo "DataScienceCluster '$DSC_NAME' exists."
else
  echo "DataScienceCluster '$DSC_NAME' does NOT exist."
  exit 1
fi
echo "Update Data Science Cluster"
kubectl patch datasciencecluster default-dsc -p '{"spec":{"components":{"modelregistry":{"managementState":"Managed"}}}}' --type=merge -o yaml
MR_NAMESPACE=$(kubectl get datasciencecluster "$DSC_NAME" -o jsonpath='{.spec.components.modelregistry.registriesNamespace}' 2>/dev/null)

echo 'Check if Namespace '$MR_NAMESPACE' exists:'
kubectl wait --for=create namespace/"$MR_NAMESPACE" --timeout="1m"
kubectl wait --for=jsonpath='{.status.phase}=Active' namespace/"$MR_NAMESPACE" --timeout="5m"
echo "Namespace '$MR_NAMESPACE' is Active."

kubectl apply -k manifests/kustomize/overlays/db-odh -n "$MR_NAMESPACE"

if ! kubectl wait --for=condition=available -n "$MR_NAMESPACE" deployment/model-registry-db --timeout=5m ; then
    kubectl events -A
    kubectl describe deployment/model-registry-db -n "$MR_NAMESPACE"
    kubectl logs deployment/model-registry-db -n "$MR_NAMESPACE"
    exit 1
fi

MYSQL_PORT_STR=$(kubectl get configmap model-registry-db-parameters -n "$MR_NAMESPACE" -o jsonpath='{.data.MYSQL_PORT}')

sed "s/DB_PORT_PLACEHOLDER/$MYSQL_PORT_STR/" "${SCRIPT_DIR}/manifests/model_registry_resource/model_registry_resource.yaml" | kubectl apply -n "$MR_NAMESPACE" -f -
kubectl wait --for=create service/model-registry -n "$MR_NAMESPACE" --timeout=5m
kubectl wait --for=condition=Available -n "$MR_NAMESPACE" deployment/model-registry --timeout=5m
kubectl port-forward -n "$MR_NAMESPACE" services/model-registry 8080:8080 & echo $! >> "${SCRIPT_DIR}/manifests/model_registry_resource/.port-forwards.pid"

