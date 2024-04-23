#!/bin/bash
FLAG_WITH_NAMESPACE=$1
MR_HOSTNAME="http://$(oc get route modelregistry-sample-http $FLAG_WITH_NAMESPACE --template='{{.spec.host}}')"
source "openshift-ci/scripts/colour_text_variables.sh"

# Function to send data and extract ID from response
make_post_extract_id() {
	local url="$1"
	local data="$2"
	local id=$(curl -s -X POST "$url" \
		-H 'accept: application/json' \
		-H 'Content-Type: application/json' \
		-d "$data" | jq -r '.id')

	if [ -z "$id" ]; then
		echo -e "${RED}Error:${NC} Error: Failed to extract ID from response"
		exit 1
  else
    echo "$id"
	fi
}

# Function to post model registry data
post_model_registry_data() {
  timestamp=$(date +"%Y%m%d%H%M%S")
  rm_name="demo-$timestamp"

  rm_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha3/registered_models" '{
    "description": "lorem ipsum registered model",
    "name": "'"$rm_name"'"
  }')

  if [ $? -ne 0 ]; then
    echo -e "${RED}Error:${NC} Registered Model ID not returned"
    exit 1
  else
    echo -e "${GREEN}✔ Success:${NC} Registered Model ID: $rm_id"
  fi

  mv_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha3/model_versions" '{
    "description": "lorem ipsum model version",
    "name": "v1",
    "author": "John Doe",
    "registeredModelId": "'"$rm_id"'"
  }')

  if [ $? -ne 0 ]; then
    echo -e "${RED}Error:${NC} Model Version ID not returned"
    exit 1
  else
     echo -e "${GREEN}✔ Success:${NC} Model Version ID: $mv_id"
  fi

  RAW_ML_MODEL_URI='https://huggingface.co/tarilabs/mnist/resolve/v1.nb20231206162408/mnist.onnx'
  ma_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha3/model_versions/$mv_id/artifacts" '{
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
    echo -e "${RED}Error:${NC} Model Artifact ID not returned"
    exit 1
  else
    echo -e "${GREEN}✔ Success:${NC} Model Artifact ID: $ma_id"
  fi
}

# Function to test deployment of second odh-project
odh_project_b() {
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
}

# Function to test Inference Service
# TODO this will continue once we have MC PR merged from: https://github.com/opendatahub-io/odh-model-controller/pull/135
inference_service() {
  iss_mr=$(curl -s -X 'GET' "$MR_HOSTNAME/api/model_registry/v1alpha3/inference_services" \
    -H 'accept: application/json')

    if [ $? -ne 0 ]; then
    echo -e "${RED}Error:${NC} InferenceService entities on MR not returned"
    exit 1
  else
    echo -e "${GREEN}✔ Success:${NC} InferenceService entities on MR: $iss_mr"
  fi
}

# Main function for orchestrating test
main() {   
  post_model_registry_data
  # odh_project_b
  # inference_service
}

# Execute main function
main
