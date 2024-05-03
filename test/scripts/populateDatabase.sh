#!/bin/bash

# This script will push data to model-registry-db. It targets the default model-registry
# To run the script from the root of the repo :  ./test/scripts/populateDatabase.sh <number of entries>
# Example :  Putting 5 entries into the database : ./test/scripts/populateDatabase.sh 5

LOOPS=$1
MR_HOSTNAME="http://$(oc get route modelregistry-sample-http -n opendatahub --template='{{.spec.host}}')"

# Function to send data and extract ID from response
make_post_extract_id() {  
    local url="$1"
    local data="$2"
    local id=$(curl -s -X POST "$url" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d "$data" | jq -r '.id')
	
  if [ -z "$id" ]; then
		echo -e "Error: Failed to extract ID from response"
		exit 1
  else
    echo "$id"
	fi
}

# Function to post model registry data
post_model_registry_data() {
  test_data_number=$1
  timestamp=$(date +"%Y%m%d%H%M%S")
  rm_name="test-data-2-$test_data_number"

  rm_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha3/registered_models" '{
    "description": "lorem ipsum registered model",
    "name": "'"$rm_name"'"
  }')

  if [ $? -ne 0 ]; then
    echo -e "Error: Registered Model ID not returned"
    exit 1
  else
    echo -e "Success: Registered Model ID: $rm_id"
  fi

  mv_id=$(make_post_extract_id "$MR_HOSTNAME/api/model_registry/v1alpha3/model_versions" '{
    "description": "lorem ipsum model version",
    "name": "v1",
    "author": "John Doe",
    "registeredModelId": "'"$rm_id"'"
  }')

  if [ $? -ne 0 ]; then
    echo -e "Error: Model Version ID not returned"
    exit 1
  else
     echo -e "Success: Model Version ID: $mv_id"
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
    echo -e "Error: Model Artifact ID not returned"
    exit 1
  else
    echo -e "Success: Model Artifact ID: $ma_id"
  fi
}

# Main function for populating database
main() {     
  for i in {$(seq 1 $LOOPS)}; do
    post_model_registry_data $i
  done
}

# Execute main function
main
