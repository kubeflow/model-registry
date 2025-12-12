#!/bin/bash

OCI_REGISTRY_NAMESPACE="local-oci-registry-ns"
set -e

echo "Undeploy local registry:"
echo "Delete Deployment"
kubectl delete deployment -l app=distribution-registry-test -n $OCI_REGISTRY_NAMESPACE
echo "Delete Service"
kubectl delete Service -l app=distribution-registry-test -n $OCI_REGISTRY_NAMESPACE
echo "Delete PersistentVolumeClaim"
kubectl delete persistentvolumeclaim distribution-registry-pvc -n $OCI_REGISTRY_NAMESPACE
echo 'Delete $OCI_REGISTRY_NAMESPACE namespace if it exists'
kubectl delete namespace $OCI_REGISTRY_NAMESPACE --wait=False

echo "Waiting for namespace $OCI_REGISTRY_NAMESPACE to be deleted..."

kubectl wait --for=delete namespace/"$OCI_REGISTRY_NAMESPACE" --timeout=300s

echo "Local registry clean up completed"
