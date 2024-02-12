#!/bin/bash

make_post_extract_id() {
    local url="$1"
    local data="$2"
    local id=$(curl -s -X POST "$url" \
        -H 'accept: application/json' \
        -H 'Content-Type: application/json' \
        -d "$data" | jq -r '.id')

    if [ -z "$id" ]; then
        echo "Error: Failed to extract ID from response"
        exit 1
    fi

    echo "$id"
}

# TODO: finalize using openshift-ci values.
OCP_CLUSTER_NAME="PROVIDE OCP CLUSTER NAME FOR OPENSHIFT-CI"
MR_NAMESPACE="shared-modelregistry-ns"
MR_HOSTNAME="http://modelregistry-sample-http-$MR_NAMESPACE.apps.$OCP_CLUSTER_NAME"

timestamp=$(date +"%Y%m%d%H%M%S")
rm_name="demo-$timestamp"

rm_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha1/registered_models" '{
  "description": "lorem ipsum registered model",
  "name": "'"$rm_name"'"
}')

if [ $? -ne 0 ]; then
    exit 1
fi
echo "Registered Model ID: $rm_id"

mv_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha1/model_versions" '{
  "description": "lorem ipsum model version",
  "name": "v1",
  "author": "John Doe",
  "registeredModelID": "'"$rm_id"'"
}')

if [ $? -ne 0 ]; then
    exit 1
fi
echo "Model Version ID: $mv_id"

RAW_ML_MODEL_URI='https://huggingface.co/tarilabs/mnist/resolve/v1.nb20231206162408/mnist.onnx'
ma_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha1/model_versions/$mv_id/artifacts" '{
  "description": "lorem ipsum model artifact",
  "uri": "'"$RAW_ML_MODEL_URI"'",
  "name": "mnist",
  "modelFormatName": "onnx",
  "modelFormatVersion": "1",
  "storageKey": "aws-connection-unused",
  "storagePath": "unused just demo",
  "artifactType": "model-artifact"
}')

if [ $? -ne 0 ]; then
    exit 1
fi
echo "Model Artifact ID: $ma_id"

ISVC_TARGET_NS=odh-project-b
MODEL_SERVER_NAME=modelserverb

oc apply -n $ISVC_TARGET_NS -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "$rm_name"
  annotations:
    "openshift.io/display-name": "$rm_name"
    "serving.kserve.io/deploymentMode": "ModelMesh"
  labels:
    "mr-registered-model-id": "$rm_id"
    "mr-model-version-id": "$mv_id"
    "mr-namespace": "$MR_NAMESPACE"
    "opendatahub.io/dashboard": "true"
spec:
  predictor:
    model:
      modelFormat:
        name: "onnx"
        version: "1"
      runtime: "$MODEL_SERVER_NAME"
      storageUri: "$RAW_ML_MODEL_URI"
EOF

# TODO this will continue once we have MC PR merged from: https://github.com/opendatahub-io/odh-model-controller/pull/135
iss_mr=$(curl -s -X 'GET' "$MR_HOSTNAME/api/model_registry/v1alpha1/inference_services" \
        -H 'accept: application/json')

echo "InferenceService entities on MR:"
echo "$iss_mr"
