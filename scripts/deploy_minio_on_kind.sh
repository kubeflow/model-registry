#!/usr/bin/env bash

set -e

DIR="$(dirname "$0")"
MINIO_NAMESPACE="minio"

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
    # echo "kubectl port-forward -n $MR_NAMESPACE service/model-registry-service 8080:8080 &"
fi

echo 'Deploying minio to Kind cluster'
if [[ $(kubectl get namespaces || false) =~ $MINIO_NAMESPACE ]]; then
    echo 'Namespace already exists, skipping creation'
else
    kubectl create namespace "$MINIO_NAMESPACE"
fi


kubectl apply -f $DIR/manifests/minio/deployment.yaml -n $MINIO_NAMESPACE
if ! kubectl wait --for=condition=available deployment/minio -n $MINIO_NAMESPACE --timeout=3m ; then 
     echo "Minio deployment took more than 3 minutes."
     kubectl events -A 
     kubectl describe deployment/minio -n $MINIO_NAMESPACE 
     kubectl logs deployment/minio -n $MINIO_NAMESPACE
     exit 1
fi

kubectl apply -f $DIR/manifests/minio/create_bucket.yaml -n $MINIO_NAMESPACE
if ! kubectl wait --for=condition=complete job/minio-init -n $MINIO_NAMESPACE --timeout=3m ; then 
     echo "Job to create Minio initialization took more than 3 minutes." 
     kubectl events -A
     kubectl logs job/minio-init -n $MINIO_NAMESPACE
     exit 1
fi

KF_MR_TEST_ACCESS_KEY_ID=$(kubectl get secret minio-secret -n minio -o jsonpath="{.data.ACCESS_KEY_ID}" | base64 --decode)
KF_MR_TEST_SECRET_ACCESS_KEY=$(kubectl get secret minio-secret -n minio -o jsonpath="{.data.SECRET_KEY}" | base64 --decode)


if [[ -z "$KF_MR_TEST_ACCESS_KEY_ID" || -z "$KF_MR_TEST_SECRET_ACCESS_KEY" ]]; then
    echo "Error: Failed to retrieve MinIO credentials. Exiting."
    exit 1
fi

cat <<EOF > $DIR/manifests/minio/.env
KF_MR_TEST_S3_ENDPOINT=http://localhost:9000
KF_MR_TEST_BUCKET_NAME=default
KF_MR_TEST_ACCESS_KEY_ID=$KF_MR_TEST_ACCESS_KEY_ID
KF_MR_TEST_SECRET_ACCESS_KEY=$KF_MR_TEST_SECRET_ACCESS_KEY
EOF
