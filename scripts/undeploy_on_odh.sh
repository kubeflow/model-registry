#!/usr/bin/env bash

set -e
SCRIPT_DIR="$(dirname "$(realpath "$BASH_SOURCE")")"
echo "Check if Data Science Cluster exists"
DSC_NAME="default-dsc"
FILE_NAME="${SCRIPT_DIR}/manifests/model_registry_resource/.port-forwards.pid"
if [ -f "$FILE_NAME" ] ; then
  echo "File $FILE_NAME exists. Deleting it..."
  while IFS= read -r pid; do
    if [[ "$pid" =~ ^[0-9]+$ ]]; then
      echo "Killing process with PID: $pid"
      kill $pid || true
    fi
  done < "$FILE_NAME"
  rm -f "$FILE_NAME"
else
  echo "File $FILE_NAME port-forward.pid does NOT exist."
fi
if kubectl get datasciencecluster "$DSC_NAME" &> /dev/null; then
  echo "DataScienceCluster '$DSC_NAME' exists."
else
  echo "DataScienceCluster '$DSC_NAME' does NOT exist."
  exit 1
fi

MR_NAMESPACE=$(kubectl get datasciencecluster "$DSC_NAME" -o jsonpath='{.spec.components.modelregistry.registriesNamespace}' 2>/dev/null)
echo "Delete modelregistry resource in namespace '$MR_NAMESPACE'."

kubectl delete modelregistry.modelregistry.opendatahub.io model-registry -n "$MR_NAMESPACE" || true
echo "Update Data Science Cluster"
kubectl patch datasciencecluster default-dsc -p '{"spec":{"components":{"modelregistry":{"managementState":"Removed"}}}}' --type=merge -o yaml

echo "Delete namespace '$MR_NAMESPACE'."
kubectl delete namespace "$MR_NAMESPACE" --wait=False
kubectl wait --for=delete namespace/"$MR_NAMESPACE" --timeout=10m
