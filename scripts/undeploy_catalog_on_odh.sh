#!/usr/bin/env bash

set -e

DSC_NAME="default-dsc"

echo "Check if DataScienceCluster exists"
if kubectl get datasciencecluster "$DSC_NAME" &> /dev/null; then
  echo "DataScienceCluster '$DSC_NAME' exists."
else
  echo "DataScienceCluster '$DSC_NAME' does NOT exist."
  exit 1
fi

echo "Check if Model Registry is enabled in DSC"
MR_STATE=$(kubectl get datasciencecluster "$DSC_NAME" -o jsonpath='{.spec.components.modelregistry.managementState}' 2>/dev/null)
if [ "$MR_STATE" != "Managed" ]; then
  echo "Model Registry is not enabled (managementState='$MR_STATE'). Expected 'Managed'."
  exit 1
fi
echo "Model Registry is enabled (managementState='Managed')."

MR_NAMESPACE=$(kubectl get datasciencecluster "$DSC_NAME" -o jsonpath='{.spec.components.modelregistry.registriesNamespace}' 2>/dev/null)
if [ -z "$MR_NAMESPACE" ]; then
  echo "Could not determine registriesNamespace from DSC."
  exit 1
fi
echo "Model Registry namespace: '$MR_NAMESPACE'"

echo "Deleting ConfigMaps to trigger operator recreation..."
for cm in mcp-catalog-sources model-catalog-sources; do
  if kubectl get configmap "$cm" -n "$MR_NAMESPACE" &> /dev/null; then
    kubectl delete configmap "$cm" -n "$MR_NAMESPACE"
    echo "ConfigMap '$cm' deleted."
  else
    echo "ConfigMap '$cm' not found, skipping."
  fi
done

echo "Restarting model-catalog pods to pick up recreated ConfigMaps"
CATALOG_DEPLOYMENT=$(kubectl get deployment -n "$MR_NAMESPACE" -l app.kubernetes.io/name=model-catalog -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
if [ -z "$CATALOG_DEPLOYMENT" ]; then
  echo "Could not find model-catalog deployment in namespace '$MR_NAMESPACE'."
  exit 1
fi
kubectl delete pod -l app.kubernetes.io/name=model-catalog -n "$MR_NAMESPACE" --wait=true
echo "Waiting for model-catalog deployment to be ready..."
kubectl wait --for=condition=Available deployment/"$CATALOG_DEPLOYMENT" -n "$MR_NAMESPACE" --timeout=5m
echo "model-catalog pod is ready."

echo "Undeploy complete."
