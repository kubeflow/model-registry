#!/bin/bash

set -e
set -o xtrace

# This test assumes there is a Kubernetes environment up and running.
# It could be either a remote one or a local one (e.g., using KinD or minikube).

DIR="$(dirname "$0")"

KUBECTL=${KUBECTL:-"kubectl"}
CLUSTER=${CLUSTER:-"kind"}

# You can provide a local version of the model registry storage initializer
# In that case, assure that is visible to the local k8s env, e.g., using 
# `kind load docker-image $MRCSI_IMG`
MRCSI_IMG=${MRCSI_IMG:-"kubeflow/model-registry-storage-initializer:main"}

KSERVE_VERSION=${KSERVE_VERSION:-"0.12"}
MODELREGISTRY_VERSION=${MODELREGISTRY_VERSION:-"v0.2.10"}
MODELREGISTRY_CSI=${MODELREGISTRY_CSI:-"v0.2.10"}

# You can provide a local model registry container image
MR_IMG=${MR_IMG:-"kubeflow/model-registry:$MODELREGISTRY_VERSION"}
# You can provide a local model registry storage initializer container image
MR_CSI_IMG=${MR_CSI_IMG:-"kubeflow/model-registry-storage-initializer:$MODELREGISTRY_CSI"}

# Check if KUBECTL is a valid command
if [ ! command -v "$KUBECTL" > /dev/null 2>&1 ]; then
    echo "KUBECTL command not found at: $KUBECTL"
    exit 1
fi

if [ ! "$KUBECTL" cluster-info > /dev/null 2>&1 ]; then
    echo "Cluster not available!"
    exit 1
fi

echo "Installing KServe version ${KSERVE_VERSION} ..."
curl -s "https://raw.githubusercontent.com/kserve/kserve/release-${KSERVE_VERSION}/hack/quick_install.sh" | bash
echo "============ KServe installed ============"

kind load docker-image -n $CLUSTER $MR_IMG
kind load docker-image -n $CLUSTER $MR_CSI_IMG

echo "Installing Model Registry ${MR_IMG} in kubeflow namespace..."
./${DIR}/../scripts/install_modelregistry.sh -i $MR_IMG
echo "======== Model Registry installed ========"

echo "Installing Model Registry ${MR_IMG} in test namespace..."
./${DIR}/../scripts/install_modelregistry.sh -i $MR_IMG -n test
echo "======== Model Registry installed ========"
