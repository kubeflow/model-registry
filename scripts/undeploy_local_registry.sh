#!/bin/bash

set -e

echo "Undeploy local registry:"
echo "Delete Deployment"
kubectl delete deployment -l app=distribution-registry-test
echo "Delete Service"
kubectl delete Service -l app=distribution-registry-test
echo "Delete PersistentVolumeClaim"
kubectl delete persistentvolumeclaim distribution-registry-pvc

echo "Local registry clean up completed"
